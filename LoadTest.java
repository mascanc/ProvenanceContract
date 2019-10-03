package com.grapevine.client;

import java.io.ByteArrayInputStream;
import java.io.FileOutputStream;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.Arrays;
import java.util.Date;
import java.util.LinkedList;
import java.util.List;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.Future;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.stream.Collectors;

import javax.xml.parsers.DocumentBuilderFactory;
import javax.xml.parsers.ParserConfigurationException;

import org.apache.commons.codec.binary.Hex;
import org.apache.xml.security.c14n.CanonicalizationException;
import org.apache.xml.security.c14n.Canonicalizer;
import org.apache.xml.security.c14n.InvalidCanonicalizerException;
import org.w3c.dom.Document;
import org.xml.sax.SAXException;

import com.grapevine.client.configuration.HLClientConfiguration;
import com.grapevine.client.configuration.HLConsumerConfiguration;
import com.grapevine.client.configuration.HLNodeConfiguration;
import com.grapevine.client.configuration.HLUserConfiguration;


/**
 * Main class to perform the load test on the blockchain.
 * @author max
 */
public class MainIJMI {

    /**
     * Data structure used to hold the information about a single
     * CDA
     * @author max
     */
    class CDAHolder {
        /** The document hash. */
        String   hash;

        /** How much it took to hash it (i.e., the canonicalisation time). */
        long     timeToHash;

        /** The CDA. */
        Document cda;

        /** The byte of the document. */
        byte[]   origDocBytes;
    }

    /**
     * The provenance client asynchronous. This is implemented the same way as the documentation
     * of Hyplerledger.
     */
    private ProvenanceClientAsync client;

    /** The list of documents to be process. */
    static List<CDAHolder>        documentList = new LinkedList<>();

    /** The list of writing time per CDA. */
    static List<Long>             writeList    = new LinkedList<>();

    /** The list of reading time per CDA. */
    static List<Long>             readList     = new LinkedList<>();

    /**
     * This function execute the async write for a single CDA.
     * @param hash
     *            The hash of provenance
     * @return The future if the execution was ok or not
     * @throws Exception
     */
    public Future<Boolean> writeSingleCda(String hash) throws Exception {

        try {
            long start = System.currentTimeMillis();
            
            // Execute with dummy data for the load test
            Future<Boolean> payload = client.insertProvenance(hash, new AgentInfo("atype", "id", "name", "idpNameId"),
                    new DoccumentLocationInfo("documentUniqueId", "id", "name", "locality"), "EX:CREATE", new Date(), null);
            long end = System.currentTimeMillis();
            long tot = end - start;
            System.out.println("Took internally " + tot);
            return payload;
        } catch (Exception e) {
            e.printStackTrace();
            throw e;
        }
    }


    /**
     * Write all the CDA of the dataset. 
     * @throws Exception
     */
    public void writeCda() throws Exception {
        long fstart = System.currentTimeMillis();

        long wstart = System.currentTimeMillis();
        int size = documentList.size();
        Future<?>[] results = new Future<?>[size];

        for (int i = 0; i < size; i++) {
            long start = System.currentTimeMillis();
            Future<Boolean> back = writeSingleCda(documentList.get(i).hash);
            results[i] = back;
            long stop = System.currentTimeMillis();
            long end = stop - start;
            System.out.println(new Date().toString() + " " + end);
            writeList.add(end);

        }
        long wend = System.currentTimeMillis();
        long wtot = wend - wstart;
        System.out.println("Write took: " + wtot);
        System.out.println("Write end. Let's now wait");
        AtomicInteger recCount = new AtomicInteger();

        /*
         * Now that we submitted, let's wait for the completion. 
         */
        long astart = System.currentTimeMillis();
        Arrays.stream(results).forEach(x -> {
            try {
                x.get();
            } catch (InterruptedException | ExecutionException e) {
                System.err.println(e.getMessage());
                recCount.incrementAndGet();
            }
        });
        long aend = System.currentTimeMillis();
        long atot = aend - astart;

        System.out.println("Wait done: " + atot + " errors " + recCount.get());

        long fend = System.currentTimeMillis();
        long ftot = fend - fstart;
        System.out.println("Took the whole writeFunction: " + ftot);
    }


    /**
     * Perform the read (sequential). 
     * 
     * @throws InvalidCanonicalizerException
     */
    private void read() throws InvalidCanonicalizerException {

        /*
         * What I do here: I get the document and I re-canonicalise it.
         */
        final Canonicalizer canon = Canonicalizer.getInstance(Canonicalizer.ALGO_ID_C14N_OMIT_COMMENTS);
        AtomicInteger readC = new AtomicInteger();
        documentList.stream().forEach(x -> {
            Document payload;
            try {
                long start = System.currentTimeMillis();
                byte[] canonBytes = canon.canonicalize(x.origDocBytes);
                payload = client.queryProvenanceDocumentByHash(new String(canonBytes));
                long end = System.currentTimeMillis();
                long tot = end - start;
                readList.add(tot);
                int mi = readC.incrementAndGet();
                if (mi % 100 == 0) {
                    System.out.println("read" + mi);

                }
            } catch (Exception e) {
                throw new RuntimeException(e);
            }
        });
    }





    /**
     * Setup the connection to the Hyperledger sample scenario on Microsoft Azure.
     * Pay attention to the firewall! The connection are tunneled via OpenSSH
     */
    public void setup() {


        try {
            HLClientConfiguration configuration = new HLClientConfiguration();
            configuration.setChannelName("masab10");
            configuration.setChainCodeName("provenanceccnew");
            configuration.setChainCodePath("ProvenanceContract");
            configuration.setChainCodeVersion("0.1");
            configuration.setTransactionWaitTime(60000);
            configuration.setUser(new HLUserConfiguration("User1", "org1.example.com", "Org1MSP",
                    "azure/crypto-config/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore/e243277a0a53eeeccdd82919a4c8e2e1737e554b71443df24b23af7a214471b8_sk",
                    "azure/crypto-config/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem"));


            configuration.setOrderer(new HLNodeConfiguration("example.com", "grpc://localhost:7050",
                    "azure/crypto-config/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem"));

            configuration.setPeer(new HLNodeConfiguration("gvwby7qwq-peer0.org1.example.com", "grpc://localhost:7051",
                    "azure/crypto-config/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem"));

            configuration.setEventHub(new HLNodeConfiguration("gvwby7qwq-peer0.org1.example.com", "grpc://localhost:7053",
                    "azure/crypto-config/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem"));

            client = new ProvenanceClientAsync(new HLConsumerConfiguration(configuration));
        } catch (Exception e) {
            e.printStackTrace();
            throw new IllegalStateException(e);
        }
    }

    /**
     * Reads all the CDA from a given path, canonicalise them and put in memory for the load test. 
     * 
     * @param path Where are the CDA stored. 
     * @throws IOException
     * @throws InvalidCanonicalizerException
     * @throws NoSuchAlgorithmException
     */
    private void readAllCDA(String path) throws IOException, InvalidCanonicalizerException, NoSuchAlgorithmException {
        org.apache.xml.security.Init.init();

        MessageDigest md = MessageDigest.getInstance("SHA-256");

        final Canonicalizer canon = Canonicalizer.getInstance(Canonicalizer.ALGO_ID_C14N_OMIT_COMMENTS);
        final AtomicInteger ai = new AtomicInteger();
        DocumentBuilderFactory dbf = DocumentBuilderFactory.newInstance();
        dbf.setNamespaceAware(true);
        documentList = Files.list(Paths.get(path)).map(x -> {

            try {
                byte[] file = Files.readAllBytes(x);
                ByteArrayInputStream bais = new ByteArrayInputStream(file);
                Document doc = dbf.newDocumentBuilder().parse(bais);



                long canonStart = System.currentTimeMillis();
                byte canonXmlBytes[] = canon.canonicalize(file);
                long canonEnd = System.currentTimeMillis();
                long tot = canonEnd - canonStart;
                String hash = Hex.encodeHexString(md.digest(canonXmlBytes));
                //   System.out.println("File " + filename + " has hash " + hash + " tot " + tot);
                ai.incrementAndGet();
                CDAHolder cda = new CDAHolder();
                cda.hash = hash;
                cda.timeToHash = tot;
                cda.cda = doc;
                cda.origDocBytes = file;
                return cda;
            } catch (SAXException | IOException | ParserConfigurationException | CanonicalizationException e) {
                throw new RuntimeException(e);
            }
        }).collect(Collectors.toList());
    }

    public static void main(String[] args) throws Exception {

        MainIJMI m = new MainIJMI();

        System.out.println("*** Setting up **** ");
        m.setup();
        /*
         * How is the code here: firstly I read all the CDAs in a memory,
         *  I canonicalise them, I evaluate the digest (and I count the timings).
         * Then I write them async, and then I read.
         */
        System.out.println("Reading all the CDAs");
        m.readAllCDA("/Users/max/Dropbox/Public/Paper Provenance/Provenance/sampledata");
        System.out.println("Read done over " + documentList.size() + " files");
        System.out.println("Now writing");
        m.writeCda();
        System.out.println("Write done");
        System.out.println("Now reading");
        m.read();
        System.out.println("Read done");

        final FileOutputStream fosCanon = new FileOutputStream("/Users/max/canontimes.txt");
        documentList.forEach(x -> {
            try {
                String toW = Long.toString(x.timeToHash);
                toW += "\n";

                fosCanon.write(toW.getBytes());
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        });
        fosCanon.close();

        final FileOutputStream fosWrite = new FileOutputStream("/Users/max/writetimes.txt");
        writeList.forEach(x -> {
            try {
                String toW = Long.toString(x);
                toW += "\n";
                fosWrite.write(toW.getBytes());

            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        });
        fosWrite.close();

        final FileOutputStream fosRead = new FileOutputStream("/Users/max/readtimes.txt");
        readList.forEach(x -> {
            try {
                String toW = Long.toString(x);
                toW += "\n";
                fosRead.write(toW.getBytes());
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        });
        fosRead.close();


    }
}
