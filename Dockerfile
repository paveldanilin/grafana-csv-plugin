FROM grafana/grafana:6.7.4-ubuntu
USER root

# Copy plugin
RUN mkdir /var/lib/grafana/plugins/grafana-csv-plugin
COPY ./dist /var/lib/grafana/plugins/grafana-csv-plugin

# Copy test data
RUN mkdir /tmp/data
COPY ./data /tmp/data
