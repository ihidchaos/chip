package dnssd

type QueryResponderSettings struct {
	mInfo *QueryResponderInfo
}

func NewQueryResponderSettings(info *QueryResponderInfo) *QueryResponderSettings {
	return &QueryResponderSettings{mInfo: info}
}

func (s *QueryResponderSettings) SetReportAdditional(qName string) *QueryResponderSettings {
	if s.IsValid() {
		s.mInfo.alsoReportAdditionalQName = true
		s.mInfo.additionalQName = qName
	}
	return s
}

func (s *QueryResponderSettings) IsValid() bool {
	return s.mInfo != nil
}

func (s *QueryResponderSettings) SetReportInServiceListing(reportService bool) *QueryResponderSettings {
	if s.IsValid() {
		s.mInfo.reportService = reportService
	}
	return s
}
