import { useEffect, useRef } from 'react';
import type { ChatMessage, Recommendation, Itinerary, BookingSummary } from '../types';
import { MessageBubble } from './MessageBubble';
import { DestinationCard } from './DestinationCard';
import { ItineraryTimeline } from './ItineraryTimeline';
import { BookingProgress } from './BookingProgress';
import { CompleteSummary } from './CompleteSummary';
import './ChatView.css';

interface ChatViewProps {
  messages: ChatMessage[];
  isStreaming: boolean;
  onSelectDestination: (id: string) => void;
  onConfirmItinerary: () => void;
  onModifyItinerary: () => void;
  onBookingPaid: (id: string) => void;
  onBookingSkip: (id: string) => void;
  onComplete: () => void;
}

export function ChatView({
  messages,
  isStreaming,
  onSelectDestination,
  onConfirmItinerary,
  onModifyItinerary,
  onBookingPaid,
  onBookingSkip,
  onComplete,
}: ChatViewProps) {
  const scrollRef = useRef<HTMLDivElement>(null);

  // Auto-scroll to bottom when new messages arrive
  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, [messages, isStreaming]);

  return (
    <div className="chat-view" ref={scrollRef} role="log" aria-live="polite" aria-label="对话消息列表">
      {messages.map((msg) => (
        <div key={msg.id} className="message-wrapper">
          <MessageBubble role={msg.role} blocks={msg.blocks} isStreaming={isStreaming} />

          {/* Render card blocks */}
          {msg.blocks.map((block, i) => {
            if (block.type !== 'card') return null;

            if (block.card_type === 'destination_recommendation') {
              const recs = block.data as Recommendation[];
              return (
                <div key={i} className="card-group">
                  {recs.map((rec) => (
                    <DestinationCard
                      key={rec.id}
                      recommendation={rec}
                      onSelect={onSelectDestination}
                    />
                  ))}
                </div>
              );
            }

            if (block.card_type === 'itinerary') {
              const itinerary = block.data as Itinerary;
              return (
                <div key={i} className="card-wrapper">
                  <ItineraryTimeline
                    itinerary={itinerary}
                    onConfirm={onConfirmItinerary}
                    onModify={onModifyItinerary}
                  />
                </div>
              );
            }

            if (block.card_type === 'booking_item') {
              // Booking items are rendered as BookingProgress (handled by booking_summary)
              return null;
            }

            if (block.card_type === 'booking_summary') {
              const summary = block.data as BookingSummary;
              if (summary.booked_count + summary.total_skipped === summary.items.length) {
                // All done, show complete summary
                return (
                  <div key={i} className="card-wrapper">
                    <CompleteSummary summary={summary} onComplete={onComplete} />
                  </div>
                );
              }
              // In progress, show booking progress
              return (
                <div key={i} className="card-wrapper">
                  <BookingProgress
                    items={summary.items}
                    onPaid={onBookingPaid}
                    onSkip={onBookingSkip}
                  />
                </div>
              );
            }

            return null;
          })}
        </div>
      ))}
    </div>
  );
}
