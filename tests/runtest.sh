curl -H "Content-Type: multipart/form-data" -X POST http://localhost:10002/ocrimage -F "lang=eng" -F"whitelist=" -F "file=@test_eng.jpg"

curl -H "Content-Type: multipart/form-data" -X POST http://localhost:10002/ocrimage -F "lang=chi_sim" -F"whitelist=" -F "file=@test_chi_sim.jpg"

curl -H "Content-Type: multipart/form-data" -X POST http://localhost:10002/ocrimage -F "lang=chi_tra" -F"whitelist=" -F "file=@test_chi_tra.jpg"
