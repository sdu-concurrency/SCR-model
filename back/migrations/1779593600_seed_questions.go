package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// Skip if the v1_multilingual record already exists
		existing, _ := app.FindFirstRecordByFilter("questions", "name = {:name}", dbx.Params{"name": "v1_multilingual"})
		if existing != nil {
			return nil
		}

		collection, err := app.FindCollectionByNameOrId("questions")
		if err != nil {
			return err
		}

		record := core.NewRecord(collection)
		record.Set("name", "supply_chain")
		record.Set("vulnerability_schema", vulnerabilitySchema)
		record.Set("capability_schema", capabilitySchema)
		record.Set("version", "v1")
		return app.Save(record)
	}, func(app core.App) error {
		// Down: remove the seeded record if it exists
		record, err := app.FindFirstRecordByFilter("questions", "name = {:name}", dbx.Params{"name": "v1_multilingual"})
		if err != nil {
			// record not found, nothing to do
			return nil
		}
		return app.Delete(record)
	})
}

// vulnerabilitySchema contains the multilingual vulnerability categories and items.
// Each item has a "value" key (e.g. "S1_1") and a "label" with "en" and "da" translations.
var vulnerabilitySchema = []map[string]any{
	{
		"label": map[string]string{"en": "1. Finance", "da": "1. Økonomi/finans"},
		"items": []map[string]any{
			{"value": "S1_1", "label": map[string]string{"en": "Asset turnover", "da": "Aktivers omsætningshastighed"}},
			{"value": "S1_2", "label": map[string]string{"en": "Access to liquidity", "da": "Adgang til likviditet"}},
			{"value": "S1_3", "label": map[string]string{"en": "High level of net working capital", "da": "Høj binding i arbejdskapital (Net Working Capital - NWC)"}},
			{"value": "S1_4", "label": map[string]string{"en": "Low cash flow", "da": "Lavt cash flow"}},
			{"value": "S1_5", "label": map[string]string{"en": "Other (Finance)", "da": "Andet (Økonomi/finans)"}},
		},
	},
	{
		"label": map[string]string{"en": "2. Customers/demand", "da": "2. Kunder/efterspørgsel"},
		"items": []map[string]any{
			{"value": "S2_1", "label": map[string]string{"en": "Unpredictability of demand", "da": "Uforudsigelig efterspørgsel"}},
			{"value": "S2_2", "label": map[string]string{"en": "Lack of sale", "da": "Manglende salg"}},
			{"value": "S2_3", "label": map[string]string{"en": "Product development pipeline", "da": "Produktudviklingspipeline"}},
			{"value": "S2_4", "label": map[string]string{"en": "Insufficient sales pipeline", "da": "Utilstrækkelig salgspipeline"}},
			{"value": "S2_5", "label": map[string]string{"en": "Customers frequently make changes in orders", "da": "Kunder ændrer hyppigt på ordrer"}},
			{"value": "S2_6", "label": map[string]string{"en": "Customer dependency", "da": "Kundeafhængighed"}},
			{"value": "S2_7", "label": map[string]string{"en": "Insufficient product assortment", "da": "Utilstrækkeligt produktsortiment"}},
			{"value": "S2_8", "label": map[string]string{"en": "Too large assortment", "da": "For stort sortiment"}},
			{"value": "S2_9", "label": map[string]string{"en": "Product liability/compensation", "da": "Produktansvar/kundeerstatning"}},
			{"value": "S2_10", "label": map[string]string{"en": "Brand image", "da": "Brandimage"}},
			{"value": "S2_11", "label": map[string]string{"en": "Time to market challenges", "da": "Time To Market udfordringer"}},
			{"value": "S2_12", "label": map[string]string{"en": "Time pressure", "da": "Presset på tid"}},
			{"value": "S2_13", "label": map[string]string{"en": "Unprofitable customers", "da": "Ikke profitable kunder"}},
			{"value": "S2_14", "label": map[string]string{"en": "Lack of market focus", "da": "Manglende markedsfokus"}},
			{"value": "S2_15", "label": map[string]string{"en": "Too low transport capacity", "da": "For lav transportkapacitet"}},
			{"value": "S2_16", "label": map[string]string{"en": "Other (Customers/demand)", "da": "Andet (Kunder/efterspørgsel)"}},
		},
	},
	{
		"label": map[string]string{"en": "3. Process/organization", "da": "3. Proces/organisation"},
		"items": []map[string]any{
			{"value": "S3_1", "label": map[string]string{"en": "Too low production capacity", "da": "For lav produktionskapacitet"}},
			{"value": "S3_2", "label": map[string]string{"en": "Reliability of equipment's", "da": "Udstyrs pålidelighed"}},
			{"value": "S3_3", "label": map[string]string{"en": "Manufacturing does not take place at the right locations", "da": "Produktion sker ikke på de rigtige lokationer"}},
			{"value": "S3_4", "label": map[string]string{"en": "Undocumented processes", "da": "Udokumenterede processer"}},
			{"value": "S3_5", "label": map[string]string{"en": "A too-high operational focus", "da": "For driftsorienteret"}},
			{"value": "S3_6", "label": map[string]string{"en": "Lack of cross-functional collaboration (silo-culture)", "da": "Mangel på tværgående samarbejde (silodannelse)"}},
			{"value": "S3_7", "label": map[string]string{"en": "Lack of human resources", "da": "Mangel på menneskelige ressourcer"}},
			{"value": "S3_8", "label": map[string]string{"en": "Lack of competencies", "da": "Mangel på kompetencer"}},
			{"value": "S3_9", "label": map[string]string{"en": "Too much tacit knowledge", "da": "For meget tavs viden"}},
			{"value": "S3_10", "label": map[string]string{"en": "Too high staff turnover", "da": "For høj personaleomsætning"}},
			{"value": "S3_11", "label": map[string]string{"en": "Too dependent on key persons", "da": "For afhængig af nøglepersoner"}},
			{"value": "S3_12", "label": map[string]string{"en": "Lack of financial resources", "da": "For få finansielle ressourcer"}},
			{"value": "S3_13", "label": map[string]string{"en": "Quality", "da": "Kvalitet"}},
			{"value": "S3_14", "label": map[string]string{"en": "Lack of maintenance", "da": "Manglende vedligehold"}},
			{"value": "S3_15", "label": map[string]string{"en": "Insufficient foundation of production (master data)", "da": "Utilstrækkeligt produktionsgrundlag"}},
			{"value": "S3_16", "label": map[string]string{"en": "Other (Process/organization)", "da": "Andet (Proces/organisation)"}},
		},
	},
	{
		"label": map[string]string{"en": "4. Systems/data", "da": "4. Systemer/data"},
		"items": []map[string]any{
			{"value": "S4_1", "label": map[string]string{"en": "Insufficient systems", "da": "Utilstrækkelige systemer"}},
			{"value": "S4_2", "label": map[string]string{"en": "Lack of IT security", "da": "Mangler IT-sikkerhed"}},
			{"value": "S4_3", "label": map[string]string{"en": "Lack of Quality Management", "da": "Mangler kvalitetsstyring"}},
			{"value": "S4_4", "label": map[string]string{"en": "Too low data quality", "da": "For lav datakvalitet"}},
			{"value": "S4_5", "label": map[string]string{"en": "Too low data accessibility", "da": "For lav datatilgængelighed"}},
			{"value": "S4_6", "label": map[string]string{"en": "Too few/wrong KPIs", "da": "For få/forkerte KPI'er"}},
			{"value": "S4_7", "label": map[string]string{"en": "Conflicting KPIs", "da": "Modstridende KPI'er"}},
			{"value": "S4_8", "label": map[string]string{"en": "Insufficient use of systems", "da": "Utilstrækkelig anvendelse af systemer"}},
			{"value": "S4_9", "label": map[string]string{"en": "Other (Systems/data)", "da": "Andet (Systemer/data)"}},
		},
	},
	{
		"label": map[string]string{"en": "5. Purchasing/sourcing", "da": "5. Indkøb/sourcing"},
		"items": []map[string]any{
			{"value": "S5_1", "label": map[string]string{"en": "Too low supplier capacity", "da": "For lav leverandørkapacitet"}},
			{"value": "S5_2", "label": map[string]string{"en": "Lack of accessibility of raw materials and supplies", "da": "Manglende tilgængelighed af råvarer og forsyninger"}},
			{"value": "S5_3", "label": map[string]string{"en": "Too low supplier reliability", "da": "For lav leverandørpålidelighed"}},
			{"value": "S5_4", "label": map[string]string{"en": "Dependency on Supplier Relations", "da": "Afhængighed af leverandørrelationer"}},
			{"value": "S5_5", "label": map[string]string{"en": "Lack of access to supplier competencies", "da": "Manglende adgang til leverandørkompetencer"}},
			{"value": "S5_6", "label": map[string]string{"en": "Too little focus on new suppliers", "da": "For lidt fokus på nye leverandører"}},
			{"value": "S5_7", "label": map[string]string{"en": "Supplier bankruptcy", "da": "Leverandørkonkurs"}},
			{"value": "S5_8", "label": map[string]string{"en": "Requirements for product purity", "da": "Krav til produktrenhed"}},
			{"value": "S5_9", "label": map[string]string{"en": "Other (Purchasing/sourcing)", "da": "Andet (Indkøb/sourcing)"}},
		},
	},
	{
		"label": map[string]string{"en": "6. Supply chain end-to-end", "da": "6. Supply chain end-to-end"},
		"items": []map[string]any{
			{"value": "S6_1", "label": map[string]string{"en": "Lack of transparency", "da": "Manglende gennemsigtighed"}},
			{"value": "S6_2", "label": map[string]string{"en": "Price pressures from customers/suppliers", "da": "Prispres fra kunder/leverandører"}},
			{"value": "S6_3", "label": map[string]string{"en": "Too high/low growth", "da": "For høj/lav vækst"}},
			{"value": "S6_4", "label": map[string]string{"en": "Import and export restrictions/channels", "da": "Im- og eksportrestriktioner/kanaler"}},
			{"value": "S6_5", "label": map[string]string{"en": "Too high complexity", "da": "For stor kompleksitet"}},
			{"value": "S6_6", "label": map[string]string{"en": "Other (Supply chain end-to-end)", "da": "Andet (Supply chain end-to-end)"}},
		},
	},
	{
		"label": map[string]string{"en": "7. Environment", "da": "7. Omgivelser"},
		"items": []map[string]any{
			{"value": "S7_1", "label": map[string]string{"en": "Geopolitical disruptions", "da": "Geopolitiske forstyrrelser"}},
			{"value": "S7_2", "label": map[string]string{"en": "Fluctuations in prices and exchange rates", "da": "Udsving i priser og valutakurser"}},
			{"value": "S7_3", "label": map[string]string{"en": "Terrorism/sabotage", "da": "Terrorisme/sabotage"}},
			{"value": "S7_4", "label": map[string]string{"en": "Espionage/theft", "da": "Spionage/tyveri"}},
			{"value": "S7_5", "label": map[string]string{"en": "Cyber-attack", "da": "Cyberangreb"}},
			{"value": "S7_6", "label": map[string]string{"en": "Competitors innovation", "da": "Konkurrenternes innovation"}},
			{"value": "S7_7", "label": map[string]string{"en": "Social/demographic/cultural changes", "da": "Sociale/demografiske/kulturelle ændringer"}},
			{"value": "S7_8", "label": map[string]string{"en": "Requirements for CSR/sustainability/ESG/UN SDGs", "da": "Krav om CSR/bæredygtighed/FN's verdensmål/ESG/kædeansvar"}},
			{"value": "S7_9", "label": map[string]string{"en": "Political regulatory changes", "da": "Politiske regulatoriske ændringer"}},
			{"value": "S7_10", "label": map[string]string{"en": "Stakeholders/NGOs", "da": "Interessenter/NGO'er"}},
			{"value": "S7_11", "label": map[string]string{"en": "Disruptions", "da": "Disruptions"}},
			{"value": "S7_12", "label": map[string]string{"en": "Unclear/lack of IPR", "da": "Uklar/manglende IP-rettigheder"}},
			{"value": "S7_13", "label": map[string]string{"en": "Strikes", "da": "Strejke(r)"}},
			{"value": "S7_14", "label": map[string]string{"en": "Other (Environment)", "da": "Andet (Omgivelser)"}},
		},
	},
}

// capabilitySchema contains the multilingual capability categories and items.
// Each item has a "value" key (e.g. "K1_1") and a "label" with "en" and "da" translations.
var capabilitySchema = []map[string]any{
	{
		"label": map[string]string{"en": "1. Finance", "da": "1. Økonomi/finans"},
		"items": []map[string]any{
			{"value": "K1_1", "label": map[string]string{"en": "Degree of asset utilization", "da": "Aktivernes udnyttelsesgrad"}},
			{"value": "K1_2", "label": map[string]string{"en": "Access to capital", "da": "Adgang til likviditet"}},
			{"value": "K1_3", "label": map[string]string{"en": "Insurances", "da": "Forsikringer"}},
			{"value": "K1_4", "label": map[string]string{"en": "Price margin", "da": "Prismargen"}},
			{"value": "K1_5", "label": map[string]string{"en": "Cost estimates/costing", "da": "For- og efterkalkulation"}},
			{"value": "K1_6", "label": map[string]string{"en": "Cash flow", "da": "Cash flow"}},
			{"value": "K1_7", "label": map[string]string{"en": "Terms of payment", "da": "Betalingsbetingelser/konditioner"}},
			{"value": "K1_8", "label": map[string]string{"en": "Focus on net working capital", "da": "Fokus på arbejdskapital (NWC)"}},
			{"value": "K1_9", "label": map[string]string{"en": "Focus on customer/product profitability", "da": "Fokus på kunde/produktlønsomhed"}},
			{"value": "K1_10", "label": map[string]string{"en": "Other (Finance)", "da": "Andet (Økonomi/finans)"}},
		},
	},
	{
		"label": map[string]string{"en": "2. Customers/demand", "da": "2. Kunder/efterspørgsel"},
		"items": []map[string]any{
			{"value": "K2_1", "label": map[string]string{"en": "Customer loyalty/retention", "da": "Kundeloyalitet/-fastholdelse"}},
			{"value": "K2_2", "label": map[string]string{"en": "Market share", "da": "Markedsandele"}},
			{"value": "K2_3", "label": map[string]string{"en": "Capability to create customer relations", "da": "Evne til at skabe kunderelationer"}},
			{"value": "K2_4", "label": map[string]string{"en": "Capability to market dispersion", "da": "Evne til markedsspredning"}},
			{"value": "K2_5", "label": map[string]string{"en": "Customer communication", "da": "Kundekommunikation"}},
			{"value": "K2_6", "label": map[string]string{"en": "Customer segmentation", "da": "Kundesegmentering"}},
			{"value": "K2_7", "label": map[string]string{"en": "Forecasting", "da": "Forecasting"}},
			{"value": "K2_8", "label": map[string]string{"en": "Collaboration on forecasting", "da": "Samarbejde om forecasting"}},
			{"value": "K2_9", "label": map[string]string{"en": "Product differentiation", "da": "Produktdifferentiering"}},
			{"value": "K2_10", "label": map[string]string{"en": "Alternative distribution channels", "da": "Alternative distributionskanaler"}},
			{"value": "K2_11", "label": map[string]string{"en": "Attractive product assortment", "da": "Attraktivt produktsortiment"}},
			{"value": "K2_12", "label": map[string]string{"en": "Sales Pipeline", "da": "Salgspipeline"}},
			{"value": "K2_13", "label": map[string]string{"en": "Development projects (and pipeline)", "da": "Udviklingsprojekter/-pipeline"}},
			{"value": "K2_14", "label": map[string]string{"en": "Faster time to market", "da": "Hurtigere Time to Market (TTM)"}},
			{"value": "K2_15", "label": map[string]string{"en": "Capability of product pruning", "da": "Evnen til at sanere"}},
			{"value": "K2_16", "label": map[string]string{"en": "Other (Customers/demand)", "da": "Andet (Kunder/efterspørgsel)"}},
		},
	},
	{
		"label": map[string]string{"en": "3. Product/processes", "da": "3. Produkt/proces"},
		"items": []map[string]any{
			{"value": "K3_1", "label": map[string]string{"en": "Component commonality", "da": "Fælles komponenter"}},
			{"value": "K3_2", "label": map[string]string{"en": "Increased standard products/components", "da": "Flere standardprodukter/komponenter"}},
			{"value": "K3_3", "label": map[string]string{"en": "Modular product design", "da": "Modulært produktdesign"}},
			{"value": "K3_4", "label": map[string]string{"en": "Capability to reduce product variability", "da": "Evne til at reducere produktvarianter"}},
			{"value": "K3_5", "label": map[string]string{"en": "Excess capacity", "da": "Reservekapacitet"}},
			{"value": "K3_6", "label": map[string]string{"en": "Manufacturing foundation", "da": "Produktionsgrundlag"}},
			{"value": "K3_7", "label": map[string]string{"en": "Phasing products in and out", "da": "Ind- og udfasning af produkter"}},
			{"value": "K3_8", "label": map[string]string{"en": "Capability to prevent errors", "da": "Evne til at forebygge fejl"}},
			{"value": "K3_9", "label": map[string]string{"en": "Quality management", "da": "Kvalitet/kvalitetsstyring"}},
			{"value": "K3_10", "label": map[string]string{"en": "Productivity/elimination of waste", "da": "Produktivitet/eliminering af spild"}},
			{"value": "K3_11", "label": map[string]string{"en": "Optimize manufacturing lead-time", "da": "Optimere produktionstid"}},
			{"value": "K3_12", "label": map[string]string{"en": "Fast Changeovers", "da": "Hurtig omstilling"}},
			{"value": "K3_13", "label": map[string]string{"en": "Flexible/scalable capacity", "da": "Fleksibel/skalerbar kapacitet"}},
			{"value": "K3_14", "label": map[string]string{"en": "Postpone manufacturing", "da": "Udskyde disponering/produktion"}},
			{"value": "K3_15", "label": map[string]string{"en": "Systematic maintenance", "da": "Systematisk vedligeholdelse"}},
			{"value": "K3_16", "label": map[string]string{"en": "Standardized workflows/processes", "da": "Standardiserede arbejdsgange/processer"}},
			{"value": "K3_17", "label": map[string]string{"en": "Documented workflows/processes", "da": "Dokumentation af arbejdsgange/processer"}},
			{"value": "K3_18", "label": map[string]string{"en": "Continuous improvement", "da": "Kontinuerlige forbedringer"}},
			{"value": "K3_19", "label": map[string]string{"en": "Manufacturing at the right locations", "da": "Produktion de rigtige steder"}},
			{"value": "K3_20", "label": map[string]string{"en": "Outsourcing - make or buy analyses", "da": "Outsourcing - Make or Buy analyser"}},
			{"value": "K3_21", "label": map[string]string{"en": "Other (Manufacturing)", "da": "Andet (Produktion)"}},
		},
	},
	{
		"label": map[string]string{"en": "4. Inventory management", "da": "4. Lagerstyring"},
		"items": []map[string]any{
			{"value": "K4_1", "label": map[string]string{"en": "IT-supported inventory management", "da": "IT-understøttet lagerstyring"}},
			{"value": "K4_2", "label": map[string]string{"en": "Location management", "da": "Lokationsstyring"}},
			{"value": "K4_3", "label": map[string]string{"en": "Safety stock", "da": "Sikkerhedslager"}},
			{"value": "K4_4", "label": map[string]string{"en": "Min/max inventory management", "da": "Min./max. styring"}},
			{"value": "K4_5", "label": map[string]string{"en": "ABC inventory management", "da": "ABC-vareanalyser"}},
			{"value": "K4_6", "label": map[string]string{"en": "Focus on death goods/obsolescence", "da": "Fokus på døde varer/ukurans"}},
			{"value": "K4_7", "label": map[string]string{"en": "Other (Inventory management)", "da": "Andet (Lagerstyring)"}},
		},
	},
	{
		"label": map[string]string{"en": "5. Sourcing/purchasing", "da": "5. Sourcing/indkøb"},
		"items": []map[string]any{
			{"value": "K5_1", "label": map[string]string{"en": "Being an attractive customer", "da": "At være en attraktiv kunde"}},
			{"value": "K5_2", "label": map[string]string{"en": "Substitution of raw materials, semi-finished products, and components", "da": "Substitution af råvarer, halvfabrikata og komponenter"}},
			{"value": "K5_3", "label": map[string]string{"en": "Increased suppliers/sources of supply", "da": "Flere leverandører/forsyningskilder"}},
			{"value": "K5_4", "label": map[string]string{"en": "Capability to create supplier relations", "da": "Evne til at skabe leverandørrelationer"}},
			{"value": "K5_5", "label": map[string]string{"en": "Prioritization (segmentation) of suppliers", "da": "Prioriterer (segmenterer) leverandører"}},
			{"value": "K5_6", "label": map[string]string{"en": "Supplier assessment and auditing", "da": "Leverandørvurdering og -audit"}},
			{"value": "K5_7", "label": map[string]string{"en": "Supplier development", "da": "Leverandørudvikling"}},
			{"value": "K5_8", "label": map[string]string{"en": "Knowledge about the supply market", "da": "Viden om forsyningsmarkedet"}},
			{"value": "K5_9", "label": map[string]string{"en": "Differentiated approach towards the suppliers", "da": "Differentieret tilgang til leverandørerne"}},
			{"value": "K5_10", "label": map[string]string{"en": "Preferred part list", "da": "Preferred part liste"}},
			{"value": "K5_11", "label": map[string]string{"en": "Preferred supplier list", "da": "Preferred leverandør liste"}},
			{"value": "K5_12", "label": map[string]string{"en": "Other (Sourcing/purchasing)", "da": "Andet (Sourcing/indkøb)"}},
		},
	},
	{
		"label": map[string]string{"en": "6. Systems and data", "da": "6. Systemer og data"},
		"items": []map[string]any{
			{"value": "K6_1", "label": map[string]string{"en": "Exchange of information - internally", "da": "Informationsudveksling – internt"}},
			{"value": "K6_2", "label": map[string]string{"en": "Exchange of information - externally", "da": "Informationsudveksling – eksternt"}},
			{"value": "K6_3", "label": map[string]string{"en": "Use of contemporary information technology", "da": "Brug af tidssvarende informationsteknologi"}},
			{"value": "K6_4", "label": map[string]string{"en": "Cyber security", "da": "Cybersikkerhed"}},
			{"value": "K6_5", "label": map[string]string{"en": "Monitoring Early Warning Signals", "da": "Overvågning af early warning signaler"}},
			{"value": "K6_6", "label": map[string]string{"en": "Ownership of master data", "da": "Ejerskab af stamdata"}},
			{"value": "K6_7", "label": map[string]string{"en": "Utilizing the potential of systems in use", "da": "Brug af systemers potentiale"}},
			{"value": "K6_8", "label": map[string]string{"en": "Other (Systems and data)", "da": "Andet (Systemer og data)"}},
		},
	},
	{
		"label": map[string]string{"en": "7. Management and organization", "da": "7. Ledelse og organisation"},
		"items": []map[string]any{
			{"value": "K7_1", "label": map[string]string{"en": "Relationship management", "da": "Relationsledelse"}},
			{"value": "K7_2", "label": map[string]string{"en": "Delegated accountability", "da": "Uddelegering og ansvarliggørelse"}},
			{"value": "K7_3", "label": map[string]string{"en": "Execution skills", "da": "Eksekveringsevne/evnen til at beslutte"}},
			{"value": "K7_4", "label": map[string]string{"en": "Employee involvement", "da": "Medarbejderinvolvering"}},
			{"value": "K7_5", "label": map[string]string{"en": "Learning/benchmarking", "da": "Læring/benchmarking"}},
			{"value": "K7_6", "label": map[string]string{"en": "Communication", "da": "Kommunikation"}},
			{"value": "K7_7", "label": map[string]string{"en": "Access to qualified labor", "da": "Adgang til kvalificeret arbejdskraft"}},
			{"value": "K7_8", "label": map[string]string{"en": "Capability to attract new employees", "da": "Evner at tiltrække nye medarbejdere"}},
			{"value": "K7_9", "label": map[string]string{"en": "Crisis management", "da": "Kriseledelse"}},
			{"value": "K7_10", "label": map[string]string{"en": "Risk management", "da": "Risikoledelse"}},
			{"value": "K7_11", "label": map[string]string{"en": "Lobbyism", "da": "Lobbyisme"}},
			{"value": "K7_12", "label": map[string]string{"en": "Width in Competency Profiles", "da": "Bredde i kompetenceprofiler"}},
			{"value": "K7_13", "label": map[string]string{"en": "Creative problem-solving", "da": "Kreativ problemløsning"}},
			{"value": "K7_14", "label": map[string]string{"en": "Focus on Core Competencies", "da": "Fokus på kernekompetencer"}},
			{"value": "K7_15", "label": map[string]string{"en": "Design for manufacturing/supply chain", "da": "Design for manufacturing/supply chain"}},
			{"value": "K7_16", "label": map[string]string{"en": "Sales and Operations Planning", "da": "Sales & Operations Planning"}},
			{"value": "K7_17", "label": map[string]string{"en": "Digitalization", "da": "Digitalisering"}},
			{"value": "K7_18", "label": map[string]string{"en": "Right KPI's functional and corporate levels", "da": "Rette KPI'er på funktions- og virksomhedsniveau"}},
			{"value": "K7_19", "label": map[string]string{"en": "Intellectual property rights", "da": "IP-rettigheder"}},
			{"value": "K7_20", "label": map[string]string{"en": "Work based on documented supply chain strategy", "da": "Arbejder ud fra nedfældet supply chain strategi"}},
			{"value": "K7_21", "label": map[string]string{"en": "Cultural understanding", "da": "kulturforståelse"}},
			{"value": "K7_22", "label": map[string]string{"en": "Cross training", "da": "Krydstræning"}},
			{"value": "K7_23", "label": map[string]string{"en": "Other (Management and organization)", "da": "Andet (Ledelse og organisation)"}},
		},
	},
}
