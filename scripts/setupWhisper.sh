
# Refer to Docs here
# https://github.com/ggml-org/whisper.cpp

cd ~
git clone https://github.com/ggerganov/whisper.cpp
sudo apt update
sudo apt install cmake -y
cd whisper.cpp
make
# Other models found here https://github.com/ggml-org/whisper.cpp/blob/master/models/README.md
./models/download-ggml-model.sh small.en