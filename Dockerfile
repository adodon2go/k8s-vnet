FROM alpine

RUN mkdir -p /etc/cnns-nse

WORKDIR /etc/cnns-nse

COPY ./nse ./nse

RUN chmod +x ./nse

ENTRYPOINT [ "./nse" ]