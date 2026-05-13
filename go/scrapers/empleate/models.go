package empleate

import "time"

type EmpleateJobList struct {
	ResponseHeader struct {
		Status int `json:"status"`
		QTime  int `json:"QTime"`
		Params struct {
			JSONWrf           string   `json:"json.wrf"`
			Facet             string   `json:"facet"`
			Sort              string   `json:"sort"`
			FacetMincount     string   `json:"facet.mincount"`
			FTopicsFacetLimit string   `json:"f.topics.facet.limit"`
			JSONNl            string   `json:"json.nl"`
			Wt                string   `json:"wt"`
			Rows              string   `json:"rows"`
			DefType           string   `json:"defType"`
			Df                string   `json:"df"`
			Fl                string   `json:"fl"`
			Q                 string   `json:"q"`
			QOp               string   `json:"q.op"`
			NAMING_FAILED     string   `json:"_"`
			FacetField        []string `json:"facet.field"`
			Fq                string   `json:"fq"`
		} `json:"params"`
	} `json:"responseHeader"`
	Response struct {
		NumFound int     `json:"numFound"`
		Start    int     `json:"start"`
		MaxScore float64 `json:"maxScore"`
		Docs     []struct {
			Entitytype              string    `json:"entitytype"`
			Jornada                 string    `json:"jornada"`
			SpeStateID              string    `json:"speStateId"`
			Cno                     string    `json:"cno"`
			Comunidad               string    `json:"comunidad,omitempty"`
			Categoria               string    `json:"categoria"`
			SpeState                string    `json:"speState"`
			EducacionF              string    `json:"educacionF"`
			EducacionS              string    `json:"educacionS"`
			Provincia               string    `json:"provincia,omitempty"`
			CheckVisible            bool      `json:"checkVisible"`
			ComunidadF              string    `json:"comunidadF,omitempty"`
			Ciudad                  string    `json:"ciudad,omitempty"`
			FechaCreacionPortal     string    `json:"fechaCreacionPortal"`
			FechaCreacionBoost      string    `json:"fechaCreacionBoost"`
			EmpresaSocial           bool      `json:"empresaSocial"`
			Contenido               string    `json:"contenido"`
			ProvinciaF              string    `json:"provinciaF,omitempty"`
			ProvinciaS              string    `json:"provinciaS,omitempty"`
			ProvinciaLimitrofe      []string  `json:"provinciaLimitrofe,omitempty"`
			ExternalID              string    `json:"externalId"`
			Titulo                  string    `json:"titulo"`
			CategoriaF              string    `json:"categoriaF"`
			CategoriaS              string    `json:"categoriaS"`
			Discapacidad            bool      `json:"discapacidad"`
			TamanoCompania2         int       `json:"tamanoCompania2"`
			PrioridadPortal         string    `json:"prioridadPortal"`
			Horario                 string    `json:"horario"`
			Contacto                string    `json:"contacto"`
			JornadaF                string    `json:"jornadaF"`
			URL                     string    `json:"url"`
			Modality                string    `json:"modality"`
			FechaCreacion           time.Time `json:"fechaCreacion"`
			ID                      string    `json:"id"`
			VisibleContact          bool      `json:"visibleContact"`
			Localizacion            string    `json:"localizacion,omitempty"`
			TrabajosOfertados       int       `json:"trabajosOfertados"`
			FechaCreacionFormateada string    `json:"fechaCreacionFormateada"`
			Pais                    string    `json:"pais"`
			Origen                  string    `json:"origen"`
			TipoContratoN           string    `json:"tipoContratoN"`
			Educacion               string    `json:"educacion"`
			TipoContrato            string    `json:"tipoContrato"`
			Cno4                    string    `json:"cno4"`
			CiudadF                 string    `json:"ciudadF,omitempty"`
			FechaRevision           time.Time `json:"fechaRevision"`
			Cno3                    string    `json:"cno3"`
			Cno2                    string    `json:"cno2"`
			PaisF                   string    `json:"paisF"`
			PaisS                   string    `json:"paisS"`
			Creador                 string    `json:"creador"`
			Score                   float64   `json:"score"`
			SubcategoriaF           string    `json:"subcategoriaF,omitempty"`
			SubcategoriaS           string    `json:"subcategoriaS,omitempty"`
			EducacionReq            string    `json:"educacionReq,omitempty"`
			EducacionReqF           string    `json:"educacionReqF,omitempty"`
			Subcategoria            string    `json:"subcategoria,omitempty"`
			NoMeInteresa            []string  `json:"noMeInteresa,omitempty"`
			SalarioMin              string    `json:"salarioMin,omitempty"`
			SalarioMax              string    `json:"salarioMax,omitempty"`
		} `json:"docs"`
	} `json:"response"`
	FacetCounts struct {
		FacetQueries struct {
		} `json:"facet_queries"`
		FacetFields struct {
			PaisF struct {
				ESPAA int `json:"ESPAÑA"`
			} `json:"paisF"`
			ProvinciaF struct {
				Madrid     int `json:"MADRID"`
				Barcelona  int `json:"BARCELONA"`
				Malaga     int `json:"MALAGA"`
				Sevilla    int `json:"SEVILLA"`
				Almeria    int `json:"ALMERIA"`
				Valencia   int `json:"VALENCIA"`
				Cantabria  int `json:"CANTABRIA"`
				Pontevedra int `json:"PONTEVEDRA"`
				Asturias   int `json:"ASTURIAS"`
				Zaragoza   int `json:"ZARAGOZA"`
				Murcia     int `json:"MURCIA"`
				ACORUA     int `json:"A CORUÑA"`
				Guipuzcoa  int `json:"GUIPUZCOA"`
				Tarragona  int `json:"TARRAGONA"`
				Vizcaya    int `json:"VIZCAYA"`
				Albacete   int `json:"ALBACETE"`
				Alicante   int `json:"ALICANTE"`
				Badajoz    int `json:"BADAJOZ"`
				Castellon  int `json:"CASTELLON"`
				Cordoba    int `json:"CORDOBA"`
				LARIOJA    int `json:"LA RIOJA"`
				Lleida     int `json:"LLEIDA"`
				Salamanca  int `json:"SALAMANCA"`
				Valladolid int `json:"VALLADOLID"`
			} `json:"provinciaF"`
			Provincia struct {
				Num12 int `json:"12"`
				Num14 int `json:"14"`
				Num15 int `json:"15"`
				Num20 int `json:"20"`
				Num25 int `json:"25"`
				Num26 int `json:"26"`
				Num28 int `json:"28"`
				Num29 int `json:"29"`
				Num30 int `json:"30"`
				Num33 int `json:"33"`
				Num36 int `json:"36"`
				Num37 int `json:"37"`
				Num39 int `json:"39"`
				Num41 int `json:"41"`
				Num43 int `json:"43"`
				Num46 int `json:"46"`
				Num47 int `json:"47"`
				Num48 int `json:"48"`
				Num50 int `json:"50"`
				Num08 int `json:"08"`
				Num04 int `json:"04"`
				Num02 int `json:"02"`
				Num03 int `json:"03"`
				Num06 int `json:"06"`
			} `json:"provincia"`
			Categoria struct {
				Num13 int `json:"13"`
				Num14 int `json:"14"`
				Num17 int `json:"17"`
				Num18 int `json:"18"`
				Num19 int `json:"19"`
				Num20 int `json:"20"`
				Num22 int `json:"22"`
				Num04 int `json:"04"`
				Num07 int `json:"07"`
				Num01 int `json:"01"`
				Num05 int `json:"05"`
				Num09 int `json:"09"`
			} `json:"categoria"`
			CategoriaF struct {
				INFORMTICATELECOMUNICACIONES      int `json:"INFORMÁTICA/TELECOMUNICACIONES"`
				INGENIERACALIDADCIENCIAS          int `json:"INGENIERÍA/CALIDAD/CIENCIAS"`
				APRENDICESPRIMEREMPLEO            int `json:"APRENDICES/PRIMER EMPLEO"`
				EDUCACINSERVICIOSSOCIALES         int `json:"EDUCACIÓN/SERVICIOS SOCIALES"`
				CUIDADOSASISTENCIAENELHOGAR       int `json:"CUIDADOS/ASISTENCIA EN EL HOGAR"`
				SALUDDEPORTE                      int `json:"SALUD/DEPORTE"`
				COMUNICACINCULTURAENTRETENIMIENTO int `json:"COMUNICACIÓN/CULTURA/ENTRETENIMIENTO"`
				METALMECNICA                      int `json:"METAL/MECÁNICA"`
				ADMINISTRACIN                     int `json:"ADMINISTRACIÓN"`
				ARQUITECTURADISEO                 int `json:"ARQUITECTURA/DISEÑO"`
				CONSTRUCCIN                       int `json:"CONSTRUCCIÓN"`
				ELECTRICIDADELECTRNICAENERGA      int `json:"ELECTRICIDAD/ELECTRÓNICA/ENERGÍA"`
			} `json:"categoriaF"`
			SubcategoriaF struct {
				INGENIERAS                         int `json:"INGENIERÍAS"`
				ANALISTASPROGRAMADORES             int `json:"ANALISTAS/PROGRAMADORES"`
				MICROINFORMTICAASISTENCIATCNICA    int `json:"MICROINFORMÁTICA/ASISTENCIA TÉCNICA"`
				SISTEMASSEGURIDADREDES             int `json:"SISTEMAS/SEGURIDAD/REDES"`
				OTRASACTIVIDADESTCNICAS            int `json:"OTRAS ACTIVIDADES TÉCNICAS"`
				PRIMEREMPLEO                       int `json:"PRIMER EMPLEO"`
				GESTINDEPROYECTOS                  int `json:"GESTIÓN DE PROYECTOS"`
				ENSEANZADEIDIOMAS                  int `json:"ENSEÑANZA DE IDIOMAS"`
				Medicina                           int `json:"MEDICINA"`
				ERPCRMBUSSINESSINTELLIGENCE        int `json:"ERP/CRM/BUSSINESS INTELLIGENCE"`
				MECNICAMANTENIMIENTO               int `json:"MECÁNICA/MANTENIMIENTO"`
				COMUNICACINPUBLICIDADMARKETING     int `json:"COMUNICACIÓN/PUBLICIDAD/MARKETING"`
				Limpieza                           int `json:"LIMPIEZA"`
				SERVICIODOMSTICO                   int `json:"SERVICIO DOMÉSTICO"`
				Telecomunicaciones                 int `json:"TELECOMUNICACIONES"`
				ARTECULTURAESPECTCULOS             int `json:"ARTE/CULTURA/ESPECTÁCULOS"`
				ASISTENCIAEINTEGRACINSOCIAL        int `json:"ASISTENCIA E INTEGRACIÓN SOCIAL"`
				BANCASEGUROS                       int `json:"BANCA/SEGUROS"`
				CUIDADODEANCIANOSNIOS              int `json:"CUIDADO DE ANCIANOS/NIÑOS"`
				ENERGASRENOVABLES                  int `json:"ENERGÍAS RENOVABLES"`
				ID                                 int `json:"I+D"`
				MEDIOAMBIENTEQUMICABIOLOGALABORATO int `json:"MEDIO AMBIENTE/QUÍMICA/BIOLOGÍA/LABORATO"`
				OTRASACTIVIDADES                   int `json:"OTRAS ACTIVIDADES"`
			} `json:"subcategoriaF"`
			Subcategoria struct {
				Num13002 int `json:"13002"`
				Num13006 int `json:"13006"`
				Num13009 int `json:"13009"`
				Num14003 int `json:"14003"`
				Num17001 int `json:"17001"`
				Num17002 int `json:"17002"`
				Num17003 int `json:"17003"`
				Num17004 int `json:"17004"`
				Num17005 int `json:"17005"`
				Num17006 int `json:"17006"`
				Num18002 int `json:"18002"`
				Num18003 int `json:"18003"`
				Num18004 int `json:"18004"`
				Num18006 int `json:"18006"`
				Num19003 int `json:"19003"`
				Num19004 int `json:"19004"`
				Num19005 int `json:"19005"`
				Num20003 int `json:"20003"`
				Num22001 int `json:"22001"`
				Num04002 int `json:"04002"`
				Num07002 int `json:"07002"`
				Num01002 int `json:"01002"`
				Num07001 int `json:"07001"`
			} `json:"subcategoria"`
			Origen struct {
				TecnoEmpleo int `json:"TECNO_EMPLEO"`
				Insertia    int `json:"INSERTIA"`
				Sne         int `json:"SNE"`
				Cogiti      int `json:"COGITI"`
				Web         int `json:"WEB"`
				Soc         int `json:"SOC"`
				Lanbide     int `json:"LANBIDE"`
				Sexpe       int `json:"SEXPE"`
			} `json:"origen"`
			TipoContrato struct {
				Indefin   int `json:"indefin"`
				Contract  int `json:"contract"`
				Determin  int `json:"determin"`
				A         int `json:"a"`
				Practic   int `json:"practic"`
				Temporal  int `json:"temporal"`
				Temporary int `json:"temporary"`
				Permanent int `json:"permanent"`
				Duracion  int `json:"duracion"`
				Especific int `json:"especific"`
				Sin       int `json:"sin"`
				Autonom   int `json:"autonom"`
				Con       int `json:"con"`
				Discapac  int `json:"discapac"`
				Format    int `json:"format"`
				Freelanc  int `json:"freelanc"`
				Person    int `json:"person"`
			} `json:"tipoContrato"`
			TipoContratoN struct {
				Num0 int `json:"0"`
				Num1 int `json:"1"`
				Num2 int `json:"2"`
				Num3 int `json:"3"`
				Num4 int `json:"4"`
				Num5 int `json:"5"`
				Num6 int `json:"6"`
			} `json:"tipoContratoN"`
			NoMeInteresa struct {
				Num1729663094 int `json:"1729663094"`
				Num1730861303 int `json:"1730861303"`
				Num1731223981 int `json:"1731223981"`
				Num1731718512 int `json:"1731718512"`
				Num1732968317 int `json:"1732968317"`
				Num1734593270 int `json:"1734593270"`
				Num1735818516 int `json:"1735818516"`
				Num1737400981 int `json:"1737400981"`
				Num1737525766 int `json:"1737525766"`
				Num1738810779 int `json:"1738810779"`
				Num1739309601 int `json:"1739309601"`
				Num1741024144 int `json:"1741024144"`
				Num1742734291 int `json:"1742734291"`
				Num1744479309 int `json:"1744479309"`
				Num1746580031 int `json:"1746580031"`
				Num1747988187 int `json:"1747988187"`
				Num1750093570 int `json:"1750093570"`
				Num1750162813 int `json:"1750162813"`
				Num1758337814 int `json:"1758337814"`
				Num1762306361 int `json:"1762306361"`
				Num1763407361 int `json:"1763407361"`
				Num1763930539 int `json:"1763930539"`
				Num1773444333 int `json:"1773444333"`
				Num1775412202 int `json:"1775412202"`
				Num1775567982 int `json:"1775567982"`
				Num1775861378 int `json:"1775861378"`
				Num1780393464 int `json:"1780393464"`
				Num1781452667 int `json:"1781452667"`
				Num1782493126 int `json:"1782493126"`
				Num1782975380 int `json:"1782975380"`
				Num1783240760 int `json:"1783240760"`
				Num1783376593 int `json:"1783376593"`
				Num1784522486 int `json:"1784522486"`
				Num1786702183 int `json:"1786702183"`
				Num1787798925 int `json:"1787798925"`
				Num1789108303 int `json:"1789108303"`
				Num1793628761 int `json:"1793628761"`
				Num1800185941 int `json:"1800185941"`
				Num1800880117 int `json:"1800880117"`
				Num1802063858 int `json:"1802063858"`
				Num1802678845 int `json:"1802678845"`
				Num1804202975 int `json:"1804202975"`
				Num1804555370 int `json:"1804555370"`
				Num1805265306 int `json:"1805265306"`
				Num1805731708 int `json:"1805731708"`
				Num1806082714 int `json:"1806082714"`
				Num1806690658 int `json:"1806690658"`
				Num1807530243 int `json:"1807530243"`
				Num1808817774 int `json:"1808817774"`
				Num1809701718 int `json:"1809701718"`
				Num1810176441 int `json:"1810176441"`
				Num1811521702 int `json:"1811521702"`
				Num1811950998 int `json:"1811950998"`
				Num1815551191 int `json:"1815551191"`
				Num1815646205 int `json:"1815646205"`
				Num1816214557 int `json:"1816214557"`
				Num1816361078 int `json:"1816361078"`
				Num1816575502 int `json:"1816575502"`
				Num1816839054 int `json:"1816839054"`
				Num1818080121 int `json:"1818080121"`
				Num1818702329 int `json:"1818702329"`
				Num1820529866 int `json:"1820529866"`
				Num1824617997 int `json:"1824617997"`
				Num1824619875 int `json:"1824619875"`
				Num1824890670 int `json:"1824890670"`
				Num1825071092 int `json:"1825071092"`
				Num1825496138 int `json:"1825496138"`
				Num1826187235 int `json:"1826187235"`
				Num1826626291 int `json:"1826626291"`
				Num1827076718 int `json:"1827076718"`
				Num1827618746 int `json:"1827618746"`
				Num1836785788 int `json:"1836785788"`
				Num1837118960 int `json:"1837118960"`
				Num1837309123 int `json:"1837309123"`
				Num1837342704 int `json:"1837342704"`
				Num1837399319 int `json:"1837399319"`
				Num1838211681 int `json:"1838211681"`
				Num1838786134 int `json:"1838786134"`
				Num1840841603 int `json:"1840841603"`
				Num1840859371 int `json:"1840859371"`
				Num1842428131 int `json:"1842428131"`
				Num1842782716 int `json:"1842782716"`
				Num1845258987 int `json:"1845258987"`
				Num1846643591 int `json:"1846643591"`
				Num1846674439 int `json:"1846674439"`
				Num1846912773 int `json:"1846912773"`
				Num1846951763 int `json:"1846951763"`
				Num1847164929 int `json:"1847164929"`
				Num1847719698 int `json:"1847719698"`
				Num1847944909 int `json:"1847944909"`
				Num1848167392 int `json:"1848167392"`
				Num1848577191 int `json:"1848577191"`
				Num1848926298 int `json:"1848926298"`
				Num1850020506 int `json:"1850020506"`
				Num1850247980 int `json:"1850247980"`
				Num1850483381 int `json:"1850483381"`
				Num1850598809 int `json:"1850598809"`
				Num1850806479 int `json:"1850806479"`
				Num1851868529 int `json:"1851868529"`
				Num1852026570 int `json:"1852026570"`
			} `json:"noMeInteresa"`
			EducacionF struct {
				OtrosTiposDeFormaciN                                                                                                                    int `json:"Otros tipos de formación"`
				ENSEANZASUNIVERSITARIASOFICIALESDEMSTER                                                                                                 int `json:"ENSEÑANZAS UNIVERSITARIAS OFICIALES DE MÁSTER"`
				SINESTUDIOS                                                                                                                             int `json:"SIN ESTUDIOS"`
				ENSEANZASDEGRADOSUPERIORDEFORMACINPROFESIONALESPECFICAYEQUIVALENTESARTESPLSTICASYDISEOYDEPORTIVAS                                       int `json:"ENSEÑANZAS DE GRADO SUPERIOR DE FORMACIÓN PROFESIONAL ESPECÍFICA Y EQUIVALENTES, ARTES PLÁSTICAS Y DISEÑO Y DEPORTIVAS."`
				ENSEANZASUNIVERSITARIASDEGRADO                                                                                                          int `json:"ENSEÑANZAS UNIVERSITARIAS DE GRADO"`
				ESTUDIOSPRIMARIOSINCOMPLETOS                                                                                                            int `json:"ESTUDIOS PRIMARIOS INCOMPLETOS"`
				ENSEANZASUNIVERSITARIASDE1ERY2CICLODESLOSEGUNDOCICLOYEQUIVALENTESLICENCIADOS                                                            int `json:"ENSEÑANZAS UNIVERSITARIAS DE 1 ER Y 2º CICLO, DE SÓLO SEGUNDO CICLO Y EQUIVALENTES (LICENCIADOS)"`
				ENSEANZASUNIVERSITARIASDEPRIMERCICLOYEQUIVALENTESOPERSONASQUEHANAPROBADO3CURSOSCOMPLETOSDEUNALICENCIATURAOCRDITOSEQUIVALENTESDIPLOMADOS int `json:"ENSEÑANZAS UNIVERSITARIAS DE PRIMER CICLO Y EQUIVALENTES O PERSONAS QUE HAN APROBADO 3 CURSOS COMPLETOS DE UNA LICENCIATURA O CRÉDITOS EQUIVALENTES (DIPLOMADOS)"`
				DOCTORADOUNIVERSITARIO                                                                                                                  int `json:"DOCTORADO UNIVERSITARIO."`
				DoctoradoOEquivalente                                                                                                                   int `json:"Doctorado o equivalente"`
				ENSEANZASDEBACHILLERATO                                                                                                                 int `json:"ENSEÑANZAS DE BACHILLERATO"`
				PROGRAMASPARALAFORMACINEINSERCINLABORALQUEPRECISANDEUNATITULACINDEESTUDIOSSECUNDARIOSDEPRIMERAETAPAPARASUREALIZACINMSDE300HORAS         int `json:"PROGRAMAS PARA LA FORMACIÓN E INSERCIÓN LABORAL QUE PRECISAN DE UNA TITULACIÓN DE ESTUDIOS SECUNDARIOS DE PRIMERA ETAPA PARA SU REALIZACIÓN (MÁS DE 300 HORAS)."`
			} `json:"educacionF"`
			FechaCreacionPortal struct {
				Two0260306T000000Z int `json:"2026-03-06T00:00:00Z"`
				Two0260302T000000Z int `json:"2026-03-02T00:00:00Z"`
				Two0260224T000000Z int `json:"2026-02-24T00:00:00Z"`
				Two0260108T000000Z int `json:"2026-01-08T00:00:00Z"`
				Two0250218T000000Z int `json:"2025-02-18T00:00:00Z"`
				Two0260220T000000Z int `json:"2026-02-20T00:00:00Z"`
				Two0260225T000000Z int `json:"2026-02-25T00:00:00Z"`
				Two0260303T000000Z int `json:"2026-03-03T00:00:00Z"`
				Two0260304T000000Z int `json:"2026-03-04T00:00:00Z"`
				Two0260305T000000Z int `json:"2026-03-05T00:00:00Z"`
				Two0260122T000000Z int `json:"2026-01-22T00:00:00Z"`
				Two0260204T000000Z int `json:"2026-02-04T00:00:00Z"`
				Two0250217T000000Z int `json:"2025-02-17T00:00:00Z"`
				Two0250225T000000Z int `json:"2025-02-25T00:00:00Z"`
				Two0260216T000000Z int `json:"2026-02-16T00:00:00Z"`
				Two0260217T000000Z int `json:"2026-02-17T00:00:00Z"`
				Two0260227T000000Z int `json:"2026-02-27T00:00:00Z"`
				Two0260107T000000Z int `json:"2026-01-07T00:00:00Z"`
				Two0260126T000000Z int `json:"2026-01-26T00:00:00Z"`
				Two0260129T000000Z int `json:"2026-01-29T00:00:00Z"`
				Two0260209T000000Z int `json:"2026-02-09T00:00:00Z"`
				Two0260226T000000Z int `json:"2026-02-26T00:00:00Z"`
				Two0260115T000000Z int `json:"2026-01-15T00:00:00Z"`
				Two0260121T000000Z int `json:"2026-01-21T00:00:00Z"`
				Two0260205T000000Z int `json:"2026-02-05T00:00:00Z"`
				Two0260218T000000Z int `json:"2026-02-18T00:00:00Z"`
				Two0260219T000000Z int `json:"2026-02-19T00:00:00Z"`
				Two0250326T000000Z int `json:"2025-03-26T00:00:00Z"`
				Two0250401T000000Z int `json:"2025-04-01T00:00:00Z"`
				Two0251222T000000Z int `json:"2025-12-22T00:00:00Z"`
				Two0260203T000000Z int `json:"2026-02-03T00:00:00Z"`
				Two0260206T000000Z int `json:"2026-02-06T00:00:00Z"`
				Two0260223T000000Z int `json:"2026-02-23T00:00:00Z"`
				Two0250210T000000Z int `json:"2025-02-10T00:00:00Z"`
				Two0250212T000000Z int `json:"2025-02-12T00:00:00Z"`
				Two0250220T000000Z int `json:"2025-02-20T00:00:00Z"`
				Two0250228T000000Z int `json:"2025-02-28T00:00:00Z"`
				Two0250310T000000Z int `json:"2025-03-10T00:00:00Z"`
				Two0250319T000000Z int `json:"2025-03-19T00:00:00Z"`
				Two0251014T000000Z int `json:"2025-10-14T00:00:00Z"`
				Two0251023T000000Z int `json:"2025-10-23T00:00:00Z"`
				Two0251216T000000Z int `json:"2025-12-16T00:00:00Z"`
				Two0251223T000000Z int `json:"2025-12-23T00:00:00Z"`
				Two0251230T000000Z int `json:"2025-12-30T00:00:00Z"`
				Two0260119T000000Z int `json:"2026-01-19T00:00:00Z"`
				Two0260120T000000Z int `json:"2026-01-20T00:00:00Z"`
				Two0260127T000000Z int `json:"2026-01-27T00:00:00Z"`
				Two0260128T000000Z int `json:"2026-01-28T00:00:00Z"`
				Two0260210T000000Z int `json:"2026-02-10T00:00:00Z"`
				Two0260211T000000Z int `json:"2026-02-11T00:00:00Z"`
				Two0260212T000000Z int `json:"2026-02-12T00:00:00Z"`
				Two0231102T000000Z int `json:"2023-11-02T00:00:00Z"`
				Two0240320T000000Z int `json:"2024-03-20T00:00:00Z"`
				Two0250203T000000Z int `json:"2025-02-03T00:00:00Z"`
				Two0250213T000000Z int `json:"2025-02-13T00:00:00Z"`
				Two0250219T000000Z int `json:"2025-02-19T00:00:00Z"`
				Two0250221T000000Z int `json:"2025-02-21T00:00:00Z"`
				Two0250224T000000Z int `json:"2025-02-24T00:00:00Z"`
				Two0250306T000000Z int `json:"2025-03-06T00:00:00Z"`
				Two0250307T000000Z int `json:"2025-03-07T00:00:00Z"`
				Two0250313T000000Z int `json:"2025-03-13T00:00:00Z"`
				Two0250320T000000Z int `json:"2025-03-20T00:00:00Z"`
				Two0250324T000000Z int `json:"2025-03-24T00:00:00Z"`
				Two0250325T000000Z int `json:"2025-03-25T00:00:00Z"`
				Two0250327T000000Z int `json:"2025-03-27T00:00:00Z"`
				Two0250328T000000Z int `json:"2025-03-28T00:00:00Z"`
				Two0250403T000000Z int `json:"2025-04-03T00:00:00Z"`
				Two0250514T000000Z int `json:"2025-05-14T00:00:00Z"`
				Two0250618T000000Z int `json:"2025-06-18T00:00:00Z"`
				Two0250707T000000Z int `json:"2025-07-07T00:00:00Z"`
				Two0251024T000000Z int `json:"2025-10-24T00:00:00Z"`
				Two0251030T000000Z int `json:"2025-10-30T00:00:00Z"`
				Two0251211T000000Z int `json:"2025-12-11T00:00:00Z"`
				Two0251217T000000Z int `json:"2025-12-17T00:00:00Z"`
				Two0251218T000000Z int `json:"2025-12-18T00:00:00Z"`
				Two0251219T000000Z int `json:"2025-12-19T00:00:00Z"`
				Two0251229T000000Z int `json:"2025-12-29T00:00:00Z"`
				Two0260112T000000Z int `json:"2026-01-12T00:00:00Z"`
				Two0260114T000000Z int `json:"2026-01-14T00:00:00Z"`
				Two0260130T000000Z int `json:"2026-01-30T00:00:00Z"`
				Two0260213T000000Z int `json:"2026-02-13T00:00:00Z"`
				Two0210630T000000Z int `json:"2021-06-30T00:00:00Z"`
				Two0211202T000000Z int `json:"2021-12-02T00:00:00Z"`
				Two0220526T000000Z int `json:"2022-05-26T00:00:00Z"`
				Two0230417T000000Z int `json:"2023-04-17T00:00:00Z"`
				Two0230531T000000Z int `json:"2023-05-31T00:00:00Z"`
				Two0240220T000000Z int `json:"2024-02-20T00:00:00Z"`
				Two0240406T000000Z int `json:"2024-04-06T00:00:00Z"`
				Two0240531T000000Z int `json:"2024-05-31T00:00:00Z"`
				Two0250205T000000Z int `json:"2025-02-05T00:00:00Z"`
				Two0250214T000000Z int `json:"2025-02-14T00:00:00Z"`
				Two0250311T000000Z int `json:"2025-03-11T00:00:00Z"`
				Two0250312T000000Z int `json:"2025-03-12T00:00:00Z"`
				Two0250318T000000Z int `json:"2025-03-18T00:00:00Z"`
				Two0250429T000000Z int `json:"2025-04-29T00:00:00Z"`
				Two0250505T000000Z int `json:"2025-05-05T00:00:00Z"`
				Two0250520T000000Z int `json:"2025-05-20T00:00:00Z"`
				Two0250523T000000Z int `json:"2025-05-23T00:00:00Z"`
				Two0250529T000000Z int `json:"2025-05-29T00:00:00Z"`
				Two0250610T000000Z int `json:"2025-06-10T00:00:00Z"`
			} `json:"fechaCreacionPortal"`
			JornadaF struct {
				Completa    int `json:"COMPLETA"`
				Parcial     int `json:"PARCIAL"`
				Indiferente int `json:"INDIFERENTE"`
				Flexible    int `json:"FLEXIBLE"`
			} `json:"jornadaF"`
			ExperienciaF struct {
				MSDEDOSAOS    int `json:"MÁS DE DOS AÑOS"`
				MSDESEISMESES int `json:"MÁS DE SEIS MESES"`
			} `json:"experienciaF"`
			Educacion struct {
				Num0  int `json:"0"`
				Num8  int `json:"8"`
				Num11 int `json:"11"`
				Num31 int `json:"31"`
				Num32 int `json:"32"`
				Num51 int `json:"51"`
				Num54 int `json:"54"`
				Num55 int `json:"55"`
				Num59 int `json:"59"`
				Num60 int `json:"60"`
				Num61 int `json:"61"`
				Num80 int `json:"80"`
			} `json:"educacion"`
			MinExperiencia struct {
				Num2   int `json:"2"`
				Num3   int `json:"3"`
				Num024 int `json:"024"`
				Num060 int `json:"060"`
			} `json:"minExperiencia"`
			Jornada struct {
				Num1 int `json:"1"`
				Num2 int `json:"2"`
				Num3 int `json:"3"`
				Num4 int `json:"4"`
			} `json:"jornada"`
			Pais struct {
				Num13  int `json:"13"`
				Num16  int `json:"16"`
				Num724 int `json:"724"`
			} `json:"pais"`
			Discapacidad struct {
				False int `json:"false"`
				True  int `json:"true"`
			} `json:"discapacidad"`
			Cno struct {
				Num1211 int `json:"1211"`
				Num1315 int `json:"1315"`
				Num1321 int `json:"1321"`
				Num2325 int `json:"2325"`
				Num2415 int `json:"2415"`
				Num2426 int `json:"2426"`
				Num2434 int `json:"2434"`
				Num2441 int `json:"2441"`
				Num2442 int `json:"2442"`
				Num2443 int `json:"2443"`
				Num2451 int `json:"2451"`
				Num2461 int `json:"2461"`
				Num2711 int `json:"2711"`
				Num2712 int `json:"2712"`
				Num2713 int `json:"2713"`
				Num2719 int `json:"2719"`
				Num2722 int `json:"2722"`
				Num2729 int `json:"2729"`
				Num3811 int `json:"3811"`
				Num3820 int `json:"3820"`
				Num4301 int `json:"4301"`
				Undef   int `json:"undef"`
			} `json:"cno"`
			Portales struct {
			} `json:"portales"`
			ShowPortalPu struct {
			} `json:"showPortalPu"`
			ShowPortalPr struct {
			} `json:"showPortalPr"`
		} `json:"facet_fields"`
		FacetDates struct {
		} `json:"facet_dates"`
		FacetRanges struct {
		} `json:"facet_ranges"`
	} `json:"facet_counts"`
}
