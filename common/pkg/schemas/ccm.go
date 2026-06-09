package schemas

// CCMOverviewOutput represents the CCM overview structured output.
type CCMOverviewOutput struct {
	TotalCost        float64                 `json:"total_cost"`
	CostTrend        float64                 `json:"cost_trend,omitempty"`
	ForecastedCost   float64                 `json:"forecasted_cost,omitempty"`
	PeriodData       []CCMPeriodCost         `json:"period_data,omitempty"`
	CloudProviders   []CCMCloudProviderCost  `json:"cloud_providers,omitempty"`
}

// CCMPeriodCost represents cost data for a time period.
type CCMPeriodCost struct {
	Period string  `json:"period"`
	Cost   float64 `json:"cost"`
}

// CCMCloudProviderCost represents cost by cloud provider.
type CCMCloudProviderCost struct {
	Provider string  `json:"provider"`
	Cost     float64 `json:"cost"`
}

// CCMCostCategoryOutput represents a cost category.
type CCMCostCategoryOutput struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// CCMCostCategoryListOutput is the structured output for list_ccm_cost_categories.
type CCMCostCategoryListOutput = PaginatedResult[CCMCostCategoryOutput]

// CCMAnomalyOutput represents a cost anomaly.
type CCMAnomalyOutput struct {
	ID               string  `json:"id"`
	AnomalyType      string  `json:"anomaly_type,omitempty"`
	ResourceName     string  `json:"resource_name,omitempty"`
	CloudProvider    string  `json:"cloud_provider,omitempty"`
	ActualCost       float64 `json:"actual_cost"`
	ExpectedCost     float64 `json:"expected_cost,omitempty"`
	AnomalyScore     float64 `json:"anomaly_score,omitempty"`
	Timestamp        string  `json:"timestamp,omitempty"`
	Status           string  `json:"status,omitempty"`
}

// CCMAnomalyListOutput is the structured output for list_ccm_anomalies.
type CCMAnomalyListOutput = PaginatedResult[CCMAnomalyOutput]

// CCMRecommendationOutput represents a cost recommendation.
type CCMRecommendationOutput struct {
	ID                  string  `json:"id"`
	RecommendationType  string  `json:"recommendation_type,omitempty"`
	ResourceName        string  `json:"resource_name,omitempty"`
	CloudProvider       string  `json:"cloud_provider,omitempty"`
	MonthlySavings      float64 `json:"monthly_savings,omitempty"`
	MonthlyCost         float64 `json:"monthly_cost,omitempty"`
	Status              string  `json:"status,omitempty"`
}

// CCMRecommendationListOutput is the structured output for list_ccm_recommendations.
type CCMRecommendationListOutput = PaginatedResult[CCMRecommendationOutput]

// CCMPerspectiveOutput represents a cost perspective.
type CCMPerspectiveOutput struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	ViewState     string `json:"view_state,omitempty"`
	ViewType      string `json:"view_type,omitempty"`
	FolderID      string `json:"folder_id,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	LastUpdatedAt string `json:"last_updated_at,omitempty"`
}

// CCMPerspectiveListOutput is the structured output for list_ccm_perspectives.
type CCMPerspectiveListOutput = PaginatedResult[CCMPerspectiveOutput]
