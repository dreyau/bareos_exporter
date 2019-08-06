package main

import (
	"github.com/dreyau/bareos_exporter/dataaccess"
	"github.com/prometheus/client_golang/prometheus"

	log "github.com/sirupsen/logrus"
)

type bareosMetrics struct {
	TotalFiles       *prometheus.Desc
	TotalBytes       *prometheus.Desc
	LastJobBytes     *prometheus.Desc
	LastJobFiles     *prometheus.Desc
	LastJobErrors    *prometheus.Desc
	LastJobTimestamp *prometheus.Desc

	LastFullJobBytes     *prometheus.Desc
	LastFullJobFiles     *prometheus.Desc
	LastFullJobErrors    *prometheus.Desc
	LastFullJobTimestamp *prometheus.Desc
}

func bareosCollector() *bareosMetrics {
	return &bareosMetrics{
		TotalFiles: prometheus.NewDesc("bareos_files_saved_total",
			"Total files saved for server during all backups for hostname combined",
			[]string{"hostname"}, nil,
		),
		TotalBytes: prometheus.NewDesc("bareos_bytes_saved_total",
			"Total bytes saved for server during all backups for hostname combined",
			[]string{"hostname"}, nil,
		),
		LastJobBytes: prometheus.NewDesc("bareos_last_backup_job_bytes_saved_total",
			"Total bytes saved during last backup for hostname",
			[]string{"hostname", "level"}, nil,
		),
		LastJobFiles: prometheus.NewDesc("bareos_last_backup_job_files_saved_total",
			"Total files saved during last backup for hostname",
			[]string{"hostname", "level"}, nil,
		),
		LastJobErrors: prometheus.NewDesc("bareos_last_backup_job_errors_occurred_while_saving_total",
			"Total errors occurred during last backup for hostname",
			[]string{"hostname", "level"}, nil,
		),
		LastJobTimestamp: prometheus.NewDesc("bareos_last_backup_job_unix_timestamp",
			"Execution timestamp of last backup for hostname",
			[]string{"hostname", "level"}, nil,
		),
		LastFullJobBytes: prometheus.NewDesc("bareos_last_full_backup_job_bytes_saved_total",
			"Total bytes saved during last full backup (Level = F) for hostname",
			[]string{"hostname"}, nil,
		),
		LastFullJobFiles: prometheus.NewDesc("bareos_last_full_backup_job_files_saved_total",
			"Total files saved during last full backup (Level = F) for hostname",
			[]string{"hostname"}, nil,
		),
		LastFullJobErrors: prometheus.NewDesc("bareos_last_full_backup_job_errors_occurred_while_saving_total",
			"Total errors occurred during last full backup (Level = F) for hostname",
			[]string{"hostname"}, nil,
		),
		LastFullJobTimestamp: prometheus.NewDesc("bareos_last_full_backup_job_unix_timestamp",
			"Execution timestamp of last full backup (Level = F) for hostname",
			[]string{"hostname"}, nil,
		),
	}
}

func (collector *bareosMetrics) Describe(ch chan<- *prometheus.Desc) {
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

func (collector *bareosMetrics) Collect(ch chan<- prometheus.Metric) {
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

		if filesErr != nil || bytesErr != nil || jobErr != nil || fullJobErr != nil {
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

		ch <- prometheus.MustNewConstMetric(collector.LastJobBytes, prometheus.CounterValue, float64(lastServerJob.JobBytes), server, lastServerJob.Level)
		ch <- prometheus.MustNewConstMetric(collector.LastJobFiles, prometheus.CounterValue, float64(lastServerJob.JobFiles), server, lastServerJob.Level)
		ch <- prometheus.MustNewConstMetric(collector.LastJobErrors, prometheus.CounterValue, float64(lastServerJob.JobErrors), server, lastServerJob.Level)
		ch <- prometheus.MustNewConstMetric(collector.LastJobTimestamp, prometheus.CounterValue, float64(lastServerJob.JobDate.Unix()), server, lastServerJob.Level)

		ch <- prometheus.MustNewConstMetric(collector.LastFullJobBytes, prometheus.CounterValue, float64(lastFullServerJob.JobBytes), server)
		ch <- prometheus.MustNewConstMetric(collector.LastFullJobFiles, prometheus.CounterValue, float64(lastFullServerJob.JobFiles), server)
		ch <- prometheus.MustNewConstMetric(collector.LastFullJobErrors, prometheus.CounterValue, float64(lastFullServerJob.JobErrors), server)
		ch <- prometheus.MustNewConstMetric(collector.LastFullJobTimestamp, prometheus.CounterValue, float64(lastFullServerJob.JobDate.Unix()), server)
	}
}
