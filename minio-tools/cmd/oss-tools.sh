filename=$1
echo $filename
curl http://minio.puhuitech.cn:9001/public/oss-tools -o oss-tools 
chmod 755 oss-tools
./oss-tools up -f $filename
