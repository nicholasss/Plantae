package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/nicholasss/plantae/internal/auth"
	"github.com/nicholasss/plantae/internal/database"
)

// === Global Variabes ===

// LangCodes is a map of their ISO 639 codes to its full language name.
// The data is sourced from Wikipedia, full sourcing information is below.
//
// Data source:
// List of ISO 639 language codes
// https://en.wikipedia.org/wiki/List_of_ISO_639_language_codes
// Date copied: 14-June-2025
// License: CC BY-SA
var LangCodes = map[string]string{
	"ab": "Abkhazian",
	"aa": "Afar",
	"af": "Afrikaans",
	"ak": "Akan",
	"sq": "Albanian",
	"am": "Amharic",
	"ar": "Arabic",
	"an": "Aragonese",
	"hy": "Armenian",
	"as": "Assamese",
	"av": "Avaric",
	"ae": "Avestan",
	"ay": "Aymara",
	"az": "Azerbaijani",
	"bm": "Bambara",
	"ba": "Bashkir",
	"eu": "Basque",
	"be": "Belarusian",
	"bn": "Bengali",
	"bi": "Bislama",
	"bs": "Bosnian",
	"br": "Breton",
	"bg": "Bulgarian",
	"my": "Burmese",
	"ca": "Catalan, Valencian",
	"ch": "Chamorro",
	"ce": "Chechen",
	"ny": "Chichewa, Chewa, Nyanja",
	"zh": "Chinese",
	"cu": "Church Slavonic, Old Slavonic, Old Church Slavonic",
	"cv": "Chuvash",
	"kw": "Cornish",
	"co": "Corsican",
	"cr": "Cree",
	"hr": "Croatian",
	"cs": "Czech",
	"da": "Danish",
	"dv": "Divehi, Dhivehi, Maldivian",
	"nl": "Dutch, Flemish",
	"dz": "Dzongkha",
	"en": "English",
	"eo": "Esperanto",
	"et": "Estonian",
	"ee": "Ewe",
	"fo": "Faroese",
	"fj": "Fijian",
	"fi": "Finnish",
	"fr": "French",
	"fy": "Western Frisian",
	"ff": "Fulah",
	"gd": "Gaelic, Scottish Gaelic",
	"gl": "Galician",
	"lg": "Ganda",
	"ka": "Georgian",
	"de": "German",
	"el": "Greek, Modern (1453–)",
	"kl": "Kalaallisut, Greenlandic",
	"gn": "Guarani",
	"gu": "Gujarati",
	"ht": "Haitian, Haitian Creole",
	"ha": "Hausa",
	"he": "Hebrew",
	"hz": "Herero",
	"hi": "Hindi",
	"ho": "Hiri Motu",
	"hu": "Hungarian",
	"is": "Icelandic",
	"io": "Ido",
	"ig": "Igbo",
	"id": "Indonesian",
	"ia": "Interlingua (International Auxiliary Language Association)",
	"ie": "Interlingue, Occidental",
	"iu": "Inuktitut",
	"ik": "Inupiaq",
	"ga": "Irish",
	"it": "Italian",
	"ja": "Japanese",
	"jv": "Javanese",
	"kn": "Kannada",
	"kr": "Kanuri",
	"ks": "Kashmiri",
	"kk": "Kazakh",
	"km": "Central Khmer",
	"ki": "Kikuyu, Gikuyu",
	"rw": "Kinyarwanda",
	"ky": "Kyrgyz, Kirghiz",
	"kv": "Komi",
	"kg": "Kongo",
	"ko": "Korean",
	"kj": "Kuanyama, Kwanyama",
	"ku": "Kurdish",
	"lo": "Lao",
	"la": "Latin",
	"lv": "Latvian",
	"li": "Limburgan, Limburger, Limburgish",
	"ln": "Lingala",
	"lt": "Lithuanian",
	"lu": "Luba-Katanga",
	"lb": "Luxembourgish, Letzeburgesch",
	"mk": "Macedonian",
	"mg": "Malagasy",
	"ms": "Malay",
	"ml": "Malayalam",
	"mt": "Maltese",
	"gv": "Manx",
	"mi": "Maori",
	"mr": "Marathi",
	"mh": "Marshallese",
	"mn": "Mongolian",
	"na": "Nauru",
	"nv": "Navajo, Navaho",
	"nd": "North Ndebele",
	"nr": "South Ndebele",
	"ng": "Ndonga",
	"ne": "Nepali",
	"no": "Norwegian",
	"nb": "Norwegian Bokmål",
	"nn": "Norwegian Nynorsk",
	"oc": "Occitan",
	"oj": "Ojibwa",
	"or": "Oriya",
	"om": "Oromo",
	"os": "Ossetian, Ossetic",
	"pi": "Pali",
	"ps": "Pashto, Pushto",
	"fa": "Persian",
	"pl": "Polish",
	"pt": "Portuguese",
	"pa": "Punjabi, Panjabi",
	"qu": "Quechua",
	"ro": "Romanian, Moldavian, Moldovan",
	"rm": "Romansh",
	"rn": "Rundi",
	"ru": "Russian",
	"se": "Northern Sami",
	"sm": "Samoan",
	"sg": "Sango",
	"sa": "Sanskrit",
	"sc": "Sardinian",
	"sr": "Serbian",
	"sn": "Shona",
	"sd": "Sindhi",
	"si": "Sinhala, Sinhalese",
	"sk": "Slovak",
	"sl": "Slovenian",
	"so": "Somali",
	"st": "Southern Sotho",
	"es": "Spanish, Castilian",
	"su": "Sundanese",
	"sw": "Swahili",
	"ss": "Swati",
	"sv": "Swedish",
	"tl": "Tagalog",
	"ty": "Tahitian",
	"tg": "Tajik",
	"ta": "Tamil",
	"tt": "Tatar",
	"te": "Telugu",
	"th": "Thai",
	"bo": "Tibetan",
	"ti": "Tigrinya",
	"to": "Tonga (Tonga Islands)",
	"ts": "Tsonga",
	"tn": "Tswana",
	"tr": "Turkish",
	"tk": "Turkmen",
	"tw": "Twi",
	"ug": "Uighur, Uyghur",
	"uk": "Ukrainian",
	"ur": "Urdu",
	"uz": "Uzbek",
	"ve": "Venda",
	"vi": "Vietnamese",
	"vo": "Volapük",
	"wa": "Walloon",
	"cy": "Welsh",
	"wo": "Wolof",
	"xh": "Xhosa",
	"ii": "Sichuan Yi, Nuosu ",
	"yi": "Yiddish",
	"yo": "Yoruba",
	"za": "Zhuang, Chuang",
	"zu": "Zulu",
}

// === Global Types ===

type apiConfig struct {
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	db                   *database.Queries
	sl                   *slog.Logger
	localAddr            string
	platform             string
	port                 string
	JWTSecret            string
	superAdminToken      string
}

// === Utilities Response Types ===

type errorResponse struct {
	Error string `json:"error"`
}

// === Utility Functions ===

// returns true if the platform is production
func platformProduction(cfg *apiConfig) bool {
	return cfg.platform == "production"
}

// returns true if the platform is not production
func platformNotProduction(cfg *apiConfig) bool {
	return cfg.platform != "production"
}

// check header for admin access token
func (cfg *apiConfig) getUserIDFromToken(r *http.Request) (uuid.UUID, error) {
	requestAccessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		return uuid.UUID{}, err
	}

	requestUserID, err := auth.ValidateJWT(requestAccessToken, cfg.JWTSecret)
	if err != nil {
		return uuid.UUID{}, err
	}

	return requestUserID, nil
}

func loadAPIConfig() (*apiConfig, error) {
	// loading vars from .env
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	// connect to log file and write out to it
	logFile, err := os.OpenFile("log/server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Unable to open log file: %q", err)
	}
	defer logFile.Close()

	logWriter := io.MultiWriter(os.Stdout, logFile)
	opts := slog.HandlerOptions{Level: slog.LevelDebug}
	sl := slog.New(slog.NewTextHandler(logWriter, &opts))

	// connecting to database
	dbURL := os.Getenv("GOOSE_DBSTRING")
	if dbURL == "" {
		return nil, err
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	dbQueries := database.New(db)
	sl.Info("Connected to database succesfully")

	// additional vars, configuration, and return
	cfg := &apiConfig{
		accessTokenDuration:  time.Hour * 2,
		refreshTokenDuration: time.Hour * 24 * 30,
		db:                   dbQueries,
		sl:                   sl,
		localAddr:            os.Getenv("LOCAL_ADDRESS"),
		platform:             os.Getenv("PLATFORM"),
		port:                 ":" + os.Getenv("PORT"),
		JWTSecret:            os.Getenv("JWT_SECRET"),
		superAdminToken:      os.Getenv("SUPER_ADMIN_TOKEN"),
	}

	// checking the config
	if cfg.localAddr == "" {
		log.Fatal("ERROR: 'LOCAL_ADDRESS' is empty, please check .env")
	}
	if cfg.platform == "" {
		log.Fatal("ERROR: 'PLATFORM' is empty, please check .env")
	} else if cfg.platform != "production" && cfg.platform != "testing" && cfg.platform != "development" {
		log.Fatal("ERROR: 'PLATFORM' is unexpected value, please check .env")
	}
	if cfg.port == "" {
		log.Fatal("ERROR: 'PORT' is empty, please check .env")
	}
	if cfg.JWTSecret == "" {
		log.Fatal("ERROR: 'JWT_SECRET' is empty, please check .env")
	}
	if cfg.superAdminToken == "" {
		log.Fatal("ERROR: 'SUPER_ADMIN_TOKEN' is empty, please check .env")
	}

	cfg.sl.Info("Config is loaded")

	return cfg, nil
}

// === Utility Response Handlers ===

func respondWithError(err error, code int, w http.ResponseWriter) {
	log.Printf("Error has occured during request: %q", err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err != nil {
		httpStatus := http.StatusText(code)
		errorResponse := fmt.Sprintf(`{"error":"%s"}`, httpStatus)
		w.Write([]byte(errorResponse))
		return
	}

	defaultError := http.StatusText(http.StatusInternalServerError)
	errorResponse := fmt.Sprintf(`{"error":"%s"}`, defaultError)
	w.Write([]byte(errorResponse))
}
