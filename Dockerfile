FROM grafana/grafana:6.7.4-ubuntu
USER root

# Grafana envs
ENV GF_LOG_LEVEL="debug"
ENV GF_LOG_MODE="file"

# Copy plugin
RUN mkdir /var/lib/grafana/plugins/grafana-csv-plugin
COPY ./dist /var/lib/grafana/plugins/grafana-csv-plugin
RUN chmod +x /var/lib/grafana/plugins/grafana-csv-plugin/grafana-csv-plugin_linux_amd64

# Copy test data
RUN mkdir /tmp/data
COPY ./data /tmp/data
RUN chown grafana:grafana /tmp/data/*.csv

# Copy showcases
RUN mkdir /tmp/showcase
COPY ./showcase /tmp/showcase
RUN chown grafana:grafana /tmp/showcase/*.json

# Provisioning datasources + dashboards
COPY ./provisioning/datasources /etc/grafana/provisioning/datasources
COPY ./provisioning/dashboards /etc/grafana/provisioning/dashboards
