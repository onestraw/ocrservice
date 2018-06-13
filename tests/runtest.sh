curl -H "Content-Type: multipart/form-data" -X POST http://localhost:10008/ocrimage -F "languages=eng" -F "whitelist=" -F "file=@test_eng.jpg"
echo '\n------------------\n'
curl -H "Content-Type: multipart/form-data" -X POST http://localhost:10008/ocrimage -F "languages=chi_sim" -F "whitelist=" -F "file=@test_chi_sim.jpg"
echo '\n------------------\n'
curl -H "Content-Type: multipart/form-data" -X POST http://localhost:10008/ocrimage -F "languages=chi_tra" -F "whitelist=" -F "file=@test_chi_tra.jpg"
