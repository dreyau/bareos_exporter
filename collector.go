package main

import (
	"bareos_exporter/dataaccess"
	"github.com/prometheus/client_golang/prometheus"

	log "github.com/sirupsen/logrus"
)

//Define a struct for you collector that contains pointers
//to prometheus descriptors for each metric you wish to expose.
//Note you can also include fields of other types if they provide utility
//but we just won't be exposing them as metrics.
type BareosMetrics struct {
	TotalFiles *prometheus.Desc
	TotalBytes *prometheus.Desc
	LastJobBytes *prometheus.Desc
	LastJobFiles *prometheus.Desc
	LastJobErrors *prometheus.Desc
	LastJobTimestamp *prometheus.Desc

	LastFullJobBytes *prometheus.Desc
	LastFullJobFiles *prometheus.Desc
	LastFullJobErrors *prometheus.Desc
	LastFullJobTimestamp *prometheus.Desc
}

func BareosCollector() *BareosMetrics {
	return &BareosMetrics{
		TotalFiles: prometheus.NewDesc("files_saved_for_hostname_total",
			"Total files saved for server during all backups combined",
			[]string{"hostname"}, nil,
		),
		TotalBytes: prometheus.NewDesc("bytes_saved_for_hostname_total",
			"Total bytes saved for server during all backups combined",
			[]string{"hostname"}, nil,
		),
		LastJobBytes: prometheus.NewDesc("last_backup_job_bytes_saved_for_hostname_total",
			"Last job that was executed for ",
			[]string{"hostname"}, nil,
		),
		LastJobFiles: prometheus.NewDesc("last_backup_job_files_saved_for_hostname_total",
			"Last job that was executed for ",
			[]string{"hostname"}, nil,
		),
		LastJobErrors: prometheus.NewDesc("last_backup_job_errors_occurred_while_saving_for_hostname_total",
			"Last job that was executed for ",
			[]string{"hostname"}, nil,
		),
		LastJobTimestamp: prometheus.NewDesc("last_backup_job_unix_timestamp_for_hostname",
			"Last job that was executed for ",
			[]string{"hostname"}, nil,
		),
		LastFullJobBytes: prometheus.NewDesc("last_full_backup_job_bytes_saved_for_hostname_total",
			"Total bytes saved by server",
			[]string{"hostname"}, nil,
		),
		LastFullJobFiles: prometheus.NewDesc("last_full_backup_job_files_saved_for_hostname_total",
			"Total bytes saved by server",
			[]string{"hostname"}, nil,
		),
		LastFullJobErrors: prometheus.NewDesc("last_full_backup_job_errors_occurred_while_saving_for_hostname_total",
			"Total bytes saved by server",
			[]string{"hostname"}, nil,
		),
		LastFullJobTimestamp: prometheus.NewDesc("last_full_backup_job_unix_timestamp_for_hostname",
			"Total bytes saved by server",
			[]string{"hostname"}, nil,
		),
	}
}

func (collector *BareosMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.TotalFiles
	ch <- collector.TotalBytes
	ch <- collector.LastJobBytes
	ch <- collector.LastJobFiles
	ch <- collector.LastJobErrors
	ch <- collector.LastJobTimestamp
	ch <- collector.LastFullJobBytes
	ch <- collector.LastFullJobFiles
	ch <- collector.LastFullJobErrors
	ch <- collector.LastFullJobTimestamp
}

func (collector *BareosMetrics) Collect(ch chan<- prometheus.Metric) {
	connection, connectionErr := dataaccess.GetConnection(connectionString)

	defer connection.DB.Close()

	if connectionErr != nil {
		log.Error(connectionErr)
		return
	}

	var servers, getServerListErr = connection.GetServerList()

	if getServerListErr != nil {
		log.Error(getServerListErr)
		return
	}

	for _, server := range servers {
		serverFiles, filesErr := connection.TotalFiles(server)
		serverBytes, bytesErr := connection.TotalBytes(server)
		lastServerJob, jobErr := connection.LastJob(server)
		lastFullServerJob, fullJobErr := connection.LastJob(server)

		if filesErr != nil || bytesErr != nil || jobErr != nil || fullJobErr != nil{
			log.Info(server)
		}

		if filesErr != nil {
			log.Error(filesErr)
		}

		if bytesErr != nil {
			log.Error(bytesErr)
		}

		if jobErr != nil {
			log.Error(jobErr)
		}

		if fullJobErr != nil {
			log.Error(fullJobErr)
		}

		ch <- prometheus.MustNewConstMetric(collector.TotalFiles, prometheus.CounterValue, float64(serverFiles.Files), server)
		ch <- prometheus.MustNewConstMetric(collector.TotalBytes, prometheus.CounterValue, float64(serverBytes.Bytes), server)

		ch <- prometheus.MustNewConstMetric(collector.LastJobBytes, prometheus.CounterValue, float64(lastServerJob.JobBytes), server)
		ch <- prometheus.MustNewConstMetric(collector.LastJobFiles, prometheus.CounterValue, float64(lastServerJob.JobFiles), server)
		ch <- prometheus.MustNewConstMetric(collector.LastJobErrors, prometheus.CounterValue, float64(lastServerJob.JobErrors), server)
		ch <- prometheus.MustNewConstMetric(collector.LastJobTimestamp, prometheus.CounterValue, float64(lastServerJob.JobDate.Unix()), server)

		ch <- prometheus.MustNewConstMetric(collector.LastFullJobBytes, prometheus.CounterValue, float64(lastFullServerJob.JobBytes), server)
		ch <- prometheus.MustNewConstMetric(collector.LastFullJobFiles, prometheus.CounterValue, float64(lastFullServerJob.JobFiles), server)
		ch <- prometheus.MustNewConstMetric(collector.LastFullJobErrors, prometheus.CounterValue, float64(lastFullServerJob.JobErrors), server)
		ch <- prometheus.MustNewConstMetric(collector.LastFullJobTimestamp, prometheus.CounterValue, float64(lastFullServerJob.JobDate.Unix()), server)
	}
}
