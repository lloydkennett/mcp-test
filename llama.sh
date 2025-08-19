#!/bin/bash

/home/xanq/mcp-test/llama.cpp/build/bin/llama-server \
    --n-gpu-layers 40 \
    --jinja \
    -fa \
    -c 40000 \
    -v \
    --log-file /home/xanq/mcp-test/llama.log \
    -hf bartowski/Qwen_Qwen3-8B-GGUF


/home/xanq/mcp-test/llama.cpp/build/bin/llama-server --jinja -fa -c 40000 -hf bartowski/Qwen_Qwen3-8B-GGUF