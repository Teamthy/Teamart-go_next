package events

// Event bus topic names and group identifiers.
const (
	MainTopic         = "events"
	RetryTopic        = "events-retry"
	DeadLetterTopic   = "events-dlq"
	DefaultGroupID    = "teamart-events"
	RetryGroupID      = "teamart-retry-processor"
	DeadLetterGroupID = "teamart-dlq-processor"
)

// CoreTopics returns the main event topics that should be defined for the platform.
var CoreTopics = []EventType{
	OrderCreated,
	PaymentCompleted,
	InventoryUpdated,
	LivestreamStarted,
	LivestreamEnded,
	ChatMessageCreated,
	ReactionSent,
	NotificationCreated,
	WalletUpdated,
	CreatorCommissionPaid,
}

// TopicNames returns the current set of core topic names.
func TopicNames() []string {
	topics := make([]string, 0, len(CoreTopics))
	for _, topic := range CoreTopics {
		topics = append(topics, string(topic))
	}
	return topics
}
