docker run --name tei-go-test-server1 --gpus '"device=0"' -p 22123:80 -d --pull always ghcr.io/huggingface/text-embeddings-inference:latest --model-id WhereIsAI/UAE-Large-V1

docker run --name tei-go-test-server2 --gpus '"device=0"' -p 22124:80 -d --pull always ghcr.io/huggingface/text-embeddings-inference:latest --model-id BAAI/bge-reranker-large --revision refs/pr/4

docker run --name tei-go-test-server3 --gpus '"device=0"' -p 22125:80 -d --pull always ghcr.io/huggingface/text-embeddings-inference:latest --model-id SamLowe/roberta-base-go_emotions

docker run --name tei-go-test-server4 --gpus '"device=0"' -p 22126:80 -d --pull always ghcr.io/huggingface/text-embeddings-inference:latest --model-id naver/efficient-splade-VI-BT-large-query --pooling splade