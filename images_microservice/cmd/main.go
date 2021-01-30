package main

import (
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/opentracing/opentracing-go"

	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/config"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/pkg/jaeger"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/pkg/logger"
)

func CheckAvatar(file multipart.File) (string, error) {
	fileHeader := make([]byte, 1024*1024*10)
	ContentType := ""
	if _, err := file.Read(fileHeader); err != nil {
		return ContentType, err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return ContentType, err
	}

	count, err := file.Seek(0, 2)
	if err != nil {
		return ContentType, err
	}
	if count > 1024*1024*10 {
		return ContentType, err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return ContentType, err
	}
	ContentType = http.DetectContentType(fileHeader)

	if ContentType != "image/jpg" && ContentType != "image/png" && ContentType != "image/jpeg" {
		return ContentType, err
	}

	return ContentType, nil
}

func main() {
	log.Println("Starting images microservice")

	configPath := config.GetConfigPath(os.Getenv("config"))
	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	log.Printf("CFG: %-v", cfg)

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s, Mode: %s",
		cfg.GRPCServer.AppVersion,
		cfg.Logger.Level,
		cfg.GRPCServer.Mode,
	)
	appLogger.Infof("Success parsed config: %#v", cfg.GRPCServer.AppVersion)

	tracer, closer, err := jaeger.InitJaeger(cfg)
	if err != nil {
		appLogger.Fatal("cannot create tracer", err)
	}
	appLogger.Info("Jaeger connected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	appLogger.Info("Opentracing connected")

	// http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
	// 	if err := r.ParseMultipartForm(1024 * 1024 * 10); err != nil {
	// 		log.Printf("ERROR: %v", err)
	// 		http.Error(w, err.Error(), 500)
	// 		return
	// 	}
	//
	// 	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024*10)
	// 	defer r.Body.Close()
	//
	// 	file, header, err := r.FormFile("avatar")
	// 	if err != nil {
	// 		log.Printf("ERROR: %v", err)
	// 		http.Error(w, err.Error(), 500)
	// 		return
	// 	}
	//
	// 	fileType, err := CheckAvatar(file)
	// 	if err != nil {
	// 		log.Printf("ERROR: %v", err)
	// 		http.Error(w, err.Error(), 500)
	// 		return
	// 	}
	// 	log.Printf("fileType: %-v", fileType)
	//
	// 	fileName := fmt.Sprintf("%s-%s", time.Now().String(), header.Filename)
	// 	log.Printf("fileName: %s", fileName)
	//
	// 	pool := sync.Pool{New: func() interface{} {
	// 		g := gift.New(
	// 			// gift.Resize(1024, 0, gift.LanczosResampling),
	// 			gift.Resize(1024, 0, gift.LanczosResampling),
	// 			gift.Contrast(20),
	// 			gift.Brightness(7),
	// 			gift.Gamma(0.5),
	// 			// gift.CropToSize(1024, 1024, gift.CenterAnchor),
	// 		)
	//
	// 		return g
	// 	}}
	//
	// 	g := pool.Get().(*gift.GIFT)
	// 	defer pool.Put(g)
	//
	// 	src, s, err := image.Decode(file)
	// 	if err != nil {
	// 		log.Printf("ERROR: %v", err)
	// 		http.Error(w, err.Error(), 500)
	// 		return
	// 	}
	// 	log.Printf("image.Decode FORMAT: %-v", s)
	//
	// 	dst := image.NewNRGBA(g.Bounds(src.Bounds()))
	// 	g.Draw(dst, src)
	//
	// 	// f, err := os.Create(header.Filename)
	// 	// if err != nil {
	// 	// 	log.Printf("ERROR: %v", err)
	// 	// 	http.Error(w, err.Error(), 500)
	// 	// 	return
	// 	// }
	// 	// defer f.Close()
	//
	// 	buf := &bytes.Buffer{}
	// 	switch fileType {
	// 	case "image/png":
	// 		err = png.Encode(buf, dst)
	// 		if err != nil {
	// 			log.Printf("ERROR: %v", err)
	// 			http.Error(w, err.Error(), 500)
	// 			break
	// 		}
	// 		log.Printf("case image/png: %s", fileType)
	// 	case "image/jpeg":
	// 		err = jpeg.Encode(buf, dst, nil)
	// 		if err != nil {
	// 			log.Printf("ERROR: %v", err)
	// 			http.Error(w, err.Error(), 500)
	// 			break
	// 		}
	// 		log.Printf("case image/png: %s", fileType)
	// 	default:
	// 		http.Error(w, "invalid image", 500)
	// 		return
	//
	// 	}
	//
	// 	w.WriteHeader(200)
	// 	w.Write(buf.Bytes())
	//
	// })
	http.ListenAndServe(":5007", nil)
}
