# OCR Service

使用开源项目搭了一个`OCR`服务的完整技术栈

- frontend 应用: UI完整拷贝自otiai10/ocrserver
- 负载均衡器: Nginx
- backend 业务层: Gin
- worker 计算层: Tesseract
- MQ 消息队列: RabbitMQ

## Architecture

                          - - - -    - - - -
                          | APP |    | APP |
                          - - - -    - - - -
                             ^          ^
                             |          |
                             |          |
                             v          v
                         - - - - - - - - - -
                         |  Loadbalancer   |
                         - - - - - - - - - -
                           ^      ^       ^
                         /        |         \
                       /          |           \
                     /            |             \
                    v             v              v
              - - - - - -     - - - - - -      - - - - - -
              | Backend |     | Backend |      | Backend |
              - - - - - -     - - - - - -      - - - - - -
                  ^              |   ^               ^ 
                  |              |   |               |
                  v              v   |               v
              - - - - - - - - - - - - - - - - - - - - - - - -
              |                 | | | |                     |
              |     RabbitMQ    |q| |q|        RPC          |
              |                 | | | |                     |
              - - - - - - - - - - - - - - - - - - - - - - - -
                  ^              |   ^            ^        ^ 
                  |              |   |            |         \
                  v              v   |            v          \
            - - - - - -      - - - - - -    - - - - - -    - - - - - -
            | Worker  |      | Worker  |    | Worker  |    | Worker  |
            - - - - - -      - - - - - -    - - - - - -    - - - - - -

## Play with docker

### build image

    hub clone onestraw/ocrservice
    cd ocrservice
    docker build -t onestraw/ocrservice .

### Play

    docker-compose up -d
    # view in browser: 127.0.0.1:10001
    cd tests/ && ./runtest.sh
    # ...
    docker-compose logs
    docker-compose down

## Todo

- 支持PDF
- 添加用户认证
- 持久化存储
- 高可用

## Reference

- [Tesseract Open Source OCR Engine](https://github.com/tesseract-ocr/tesseract)
- [Simple OCR server](https://github.com/otiai10/ocrserver)
- [RabbitMQ RPC](http://www.rabbitmq.com/tutorials/tutorial-six-go.html)
