FROM docker.elastic.co/logstash/logstash:7.16.2
LABEL maintainer Tomoyuki Sugiyama <tomoyuki.sugiyama.j7u@jpgr.denso.com>

COPY ./config /usr/share/logstash/config
COPY ./pipeline /usr/share/logstash/pipeline
COPY ./extra_patterns/date_jp /opt/logstash/extra_patterns
