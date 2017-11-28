package main

import (
	"github.com/beevik/etree"
	"log"
	"os"
	"time"
)

// createWasGeneratedBy creates the relation xml fragment
func createWasGeneratedBy(doc *etree.Document, genTime string) (*etree.Element, error) {
	provWasGeneratedBy := doc.CreateElement("prov:wasGeneratedBy")
	provWasGenEntity := doc.CreateElement("prov:entity")
	provWasGeneratedBy.AddChild(provWasGenEntity)
	provWasGenEntity.CreateAttr("xmlns:ns1", "http://www.w3.org/ns/prov#")
	provWasGenEntity.CreateAttr("ns1:ref", "theobject")
	provWasGenActivity := doc.CreateElement("prov:activity")
	provWasGeneratedBy.AddChild(provWasGenActivity)
	provWasGenActivity.CreateAttr("xmlns:ns1", "http://www.w3.org/ns/prov#")
	provWasGenActivity.CreateAttr("ns1:ref", "theobjectcreation")
	provWasGenTime := doc.CreateElement("prov:time")
	provWasGeneratedBy.AddChild(provWasGenTime)
	_,errTime := time.Parse("2006-01-02T15:04:05.000Z", genTime)
	if errTime != nil{
		return nil,errTime
	} 
	
	provWasGenTime.SetText(genTime)
	return provWasGeneratedBy, nil
}

// createWasAssociated creates the wasAssociatedWith xml fragment
func createWasAssociated(doc *etree.Document, agentInfo agent) *etree.Element {
	provWasAssociatedWith := doc.CreateElement("prov:wasAssociatedWith")
	provWasAssactivity := doc.CreateElement("prov:activity")
	provWasAssociatedWith.AddChild(provWasAssactivity)
	provWasAssactivity.CreateAttr("xmlns:ns1", "http://www.w3.org/ns/prov#")
	provWasAssactivity.CreateAttr("ns1:ref", "theobjectcreation")
	provWasAssagent := doc.CreateElement("prov:agent")
	provWasAssociatedWith.AddChild(provWasAssagent)
	provWasAssagent.CreateAttr("xmlns:ns1", "http://www.w3.org/ns/prov#")
	provWasAssagent.CreateAttr("ns1:ref", agentInfo.id)
	return provWasAssociatedWith
}

// createEntity creates the Entity xml fragment
func createEntity(doc *etree.Document, myHash, description string) *etree.Element {
	proventity := doc.CreateElement("prov:entity")
	proventity.CreateAttr("xmlns:ns1", "http://www.w3.org/ns/prov#")
	proventity.CreateAttr("ns1:id", "theobject")

	entitylabel := doc.CreateElement("prov:label")
	entitylabel.SetText(description)
	proventity.AddChild(entitylabel)

	entitylocationlabel := doc.CreateElement("prov:location")
	proventity.AddChild(entitylocationlabel)

	entitytypelabel := doc.CreateElement("prov:type")
	entitytypelabel.SetText("XML")
	proventity.AddChild(entitytypelabel)

	entityvalue := doc.CreateElement("prov:value")
	entityvalue.SetText(myHash)
	proventity.AddChild(entityvalue)
	return proventity
}

// createDocument create the main PROV document
func createDocument(doc *etree.Document) *etree.Element {
	provdocument := doc.CreateElement("prov:document")
	provdocument.CreateAttr("xmlns:prov", "http://www.w3.org/ns/prov#")
	provdocument.CreateAttr("xmlns:xsi", "http://www.w3.org/2001/XMLSchema-instance")
	provdocument.CreateAttr("xmlns:ex", "urn:tiani:prova")
	return provdocument
}

// createActivity create the activity xml fragment
func createActivity(doc *etree.Document, activity string) *etree.Element {
	provactivity := doc.CreateElement("prov:activity")
	provactivity.CreateAttr("xmlns:ns1", "http://www.w3.org/ns/prov#")
	provactivity.CreateAttr("ns1:id", "theobjectcreation")

	activitytype := doc.CreateElement("prov:type")
	activitytype.SetText(activity)
	provactivity.AddChild(activitytype)
	return provactivity

}

// createAgent creates the Agent xml fragment
func createAgent(doc *etree.Document, agentInfo agent) *etree.Element {
	provagent := doc.CreateElement("prov:agent")
	provagent.CreateAttr("xmlns:ns1", "http://www.w3.org/ns/prov#")
	provagent.CreateAttr("ns1:id", agentInfo.id)

	agenttype := doc.CreateElement("prov:type")
	agenttype.SetText(agentInfo.atype)
	provagent.AddChild(agenttype)

	agentDocId := doc.CreateElement("hpd:doctorid")
	provagent.AddChild(agentDocId)
	agentDocId.CreateAttr("xmlns:hpd", "IHEHPD")
	agentDocId.SetText(agentInfo.id)

	agentDocName := doc.CreateElement("hpd:doctorname")
	provagent.AddChild(agentDocName)
	agentDocName.CreateAttr("xmlns:hpd", "IHEHPD")
	agentDocName.SetText(agentInfo.name)

	agentDocIdp := doc.CreateElement("hpd:idp")
	provagent.AddChild(agentDocIdp)
	agentDocIdp.CreateAttr("xmlns:hpd", "idp")
	agentDocIdp.SetText(agentInfo.identityProvider)
	return provagent
}

// makeProvenanceDocumentSegmented creates the object/xml provenance structure,
// which is for the segmented data
func makeProvenanceDocumentSegmented(hashOfTheSegment string, myHash string, agentInfo agent, activity, genTime string) (mapDocument *etree.Document, document string, err error) {
	log.Printf("creating a provenance segmented document")
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)

	provdocument := createDocument(doc)

	
	// This is the definition of the entity: here I add a generic "the cda" and
	// it shouldn't bother. This is the originating CDA
	
	proventity := createEntity(doc, myHash, "The object document")
	provdocument.AddChild(proventity)


	// now I create another entity with the hash
	provSegEntity := createEntity(doc, hashOfTheSegment, "The CDA Segment")
	provdocument.AddChild(provSegEntity)


	// this is the activity
	provactivity := createActivity(doc, activity)
	provdocument.AddChild(provactivity)

	
	// this is the agent
	provagent := createAgent(doc, agentInfo)
	provdocument.AddChild(provagent)
	
	// Now we need to say what happened. Basically here we create
	// that the entity has been generated by the activity, whose activity was
	// associated with an agent, and the eneity was attributed to the agent

	provWasGeneratedBy, errTime := createWasGeneratedBy(doc, genTime)
	if errTime != nil{
		return nil,"",errTime
	} 
	provdocument.AddChild(provWasGeneratedBy)

	// Now, who created the docuemnt, firstly we know that the agent perfomed an
	// activity, than that the entity was attributed to the agent

	provWasAssociatedWith := createWasAssociated(doc, agentInfo)
	provdocument.AddChild(provWasAssociatedWith)

	// Used
	provUsed := doc.CreateElement("prov:used")
	provdocument.AddChild(provUsed)
	provUsedactivity := doc.CreateElement("prov:activity")
	provUsed.AddChild(provUsedactivity)
	provUsedactivity.CreateAttr("xmlns:ns1", "http://www.w3.org/ns/prov#")
	provUsedactivity.CreateAttr("ns1:ref", "theobjectcreation")
	provUsedagent := doc.CreateElement("prov:entity")
	provUsed.AddChild(provUsedagent)
	provUsedagent.CreateAttr("xmlns:ns1", "http://www.w3.org/ns/prov#")
	provUsedagent.CreateAttr("ns1:ref", "theobject")

	// Was derivedFrom
	provWasDerivedFrom := doc.CreateElement("prov:wasDerivedFrom")
	provdocument.AddChild(provWasDerivedFrom)
	provUsedGeneratedEntity := doc.CreateElement("prov:generatedEntity")
	provWasDerivedFrom.AddChild(provUsedGeneratedEntity)
	provUsedGeneratedEntity.CreateAttr("xmlns:ns1", "http://www.w3.org/ns/prov#")
	provUsedGeneratedEntity.CreateAttr("ns1:ref", "thesegment")
	provUsedEntity := doc.CreateElement("prov:usedEntity")
	provWasDerivedFrom.AddChild(provUsedEntity)
	provUsedEntity.CreateAttr("xmlns:ns1", "http://www.w3.org/ns/prov#")
	provUsedEntity.CreateAttr("ns1:ref", "theobject")
	mydocument, err := doc.WriteToString()
	if err != nil {
		return nil,"",err
	}
	doc.WriteTo(os.Stdout)
	return doc, mydocument,nil
}


// makeProvenanceDocument creates the object/xml provenance structure
func makeProvenanceDocument(myHash string, agentInfo agent, activity, genTime string) (provdoc *etree.Document, docAsStr string, err error) {
	log.Printf("creating a provenance document")
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)

	provdocument := createDocument(doc)

	//
	// This is the definition of the entity: here I add a generic "the cda" and
	// it shouldn't bother
	proventity := createEntity(doc, myHash,"The object document")
	provdocument.AddChild(proventity)


	// this is the activity
	provactivity := createActivity(doc, activity)
	provdocument.AddChild(provactivity)

	// this is the agent
	provagent := createAgent(doc, agentInfo)
	provdocument.AddChild(provagent)

	// Now we need to say what happened. Basically here we create
	// that the entity has been generated by the activity, whose activity was
	// associated with an agent, and the eneity was attributed to the agent

	provWasGeneratedBy, errTime := createWasGeneratedBy(doc, genTime)
	if errTime != nil{
		return nil,"",errTime
	} 
	provdocument.AddChild(provWasGeneratedBy)

	// Now, who created the docuemnt, firstly we know that the agent perfomed an
	// activity, than that the entity was attributed to the agent

	provWasAssociatedWith := createWasAssociated(doc, agentInfo)
	provdocument.AddChild(provWasAssociatedWith)

	provWasAttributedTo := doc.CreateElement("prov:wasAttributedTo")
	provdocument.AddChild(provWasAttributedTo)
	provWasAttributedToEntity := doc.CreateElement("prov:entity")
	provWasAttributedTo.AddChild(provWasAttributedToEntity)
	provWasAttributedToEntity.CreateAttr("xmlns:ns1", "http://www.w3.org/ns/prov#")
	provWasAttributedToEntity.CreateAttr("ns1:ref", "theobject")
	provWasAttributeToagent := doc.CreateElement("prov:agent")
	provWasAttributedTo.AddChild(provWasAttributeToagent)
	provWasAttributeToagent.CreateAttr("xmlns:ns1", "http://www.w3.org/ns/prov#")
	provWasAttributeToagent.CreateAttr("ns1:ref", agentInfo.id)

	docAsString, err := doc.WriteToString()
	if err != nil {
		return nil,"",err
	}
	doc.WriteTo(os.Stdout)

	return doc, docAsString,nil
}
