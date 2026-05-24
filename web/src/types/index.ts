// ===== Enums =====

export type Phase = 'COLLECTING' | 'RECOMMENDING' | 'PLANNING' | 'BOOKING' | 'COMPLETED';

export type BookingType = 'OUTBOUND' | 'HOTEL' | 'RETURN' | 'TICKET';

export type BookingStatus = 'PENDING' | 'BOOKED' | 'SKIPPED';

export type CardType =
  | 'destination_recommendation'
  | 'itinerary'
  | 'booking_item'
  | 'booking_summary';

// ===== SSE Events =====

export interface TextEventData {
  content: string;
  delta: boolean;
}

export interface ThinkingEventData {
  content: string;
}

export interface ToolCallEventData {
  name: string;
  args: Record<string, unknown>;
  status: 'running' | 'done' | 'error';
  result: unknown | null;
}

export interface CardEventData {
  card_type: CardType;
  data: unknown;
}

export interface PhaseChangeEventData {
  from: Phase;
  to: Phase;
}

export interface ErrorEventData {
  code: string;
  message: string;
}

// ===== Business Data =====

export interface Session {
  id: string;
  phase: Phase;
  created_at: string;
  updated_at: string;
  user_request: UserRequest | null;
  recommendations: Recommendation[];
  itinerary: Itinerary | null;
  booking_items: BookingItem[];
  chat_history: ChatMessage[];
}

export interface UserRequest {
  raw_input: string;
  departure_city: string;
  travel_start: string;
  travel_end: string;
  budget: number;
  preferences: string[];
  travelers: number;
}

export interface Recommendation {
  id: string;
  destination: string;
  reason: string;
  weather: WeatherInfo;
  transport: TransportInfo;
  hotel_range: HotelRange;
  total_estimate: CostRange;
  sample_tickets: TicketInfo[];
  is_mock: boolean;
}

export interface WeatherInfo {
  temp_range: string;
  condition: string;
  rain_prob: number;
  is_mock: boolean;
}

export interface TransportInfo {
  mode: 'flight' | 'train';
  duration: string;
  price: number;
  is_mock: boolean;
}

export interface HotelRange {
  economy: number;
  comfort: number;
  is_mock: boolean;
}

export interface CostRange {
  min: number;
  max: number;
}

export interface TicketInfo {
  type: 'flight' | 'train';
  number: string;
  depart: string;
  arrive: string;
  price: number;
  seat: string;
  is_mock: boolean;
}

export interface Itinerary {
  destination: string;
  outbound: FlightDetail;
  hotel: HotelDetail;
  daily_plans: DayPlan[];
  return: FlightDetail;
  cost_summary: CostSummary;
  total_cost: number;
}

export interface FlightDetail {
  type: 'flight' | 'train';
  number: string;
  depart_city: string;
  arrive_city: string;
  depart_time: string;
  arrive_time: string;
  duration: string;
  price: number;
  seat: string;
  is_mock: boolean;
}

export interface HotelDetail {
  name: string;
  location: string;
  room_type: string;
  check_in: string;
  check_out: string;
  nights: number;
  price_per_night: number;
  total_price: number;
  is_mock: boolean;
}

export interface DayPlan {
  date: string;
  day_number: number;
  title: string;
  activities: Activity[];
}

export interface Activity {
  time: string;
  name: string;
  transport: string;
  cost: number;
  note: string;
}

export interface CostSummary {
  outbound_cost: number;
  return_cost: number;
  hotel_cost: number;
  activity_cost: number;
  meal_cost: number;
  transport_cost: number;
}

export interface BookingItem {
  id: string;
  item_type: BookingType;
  item_name: string;
  detail: string;
  platform: string;
  link: string;
  alt_links: PlatformLink[];
  price: number;
  status: BookingStatus;
  actual_cost: number;
  order: number;
}

export interface PlatformLink {
  platform: string;
  link: string;
}

export interface BookingSummary {
  destination: string;
  items: BookingItem[];
  total_booked: number;
  total_skipped: number;
  booked_count: number;
}

// ===== Message Model =====

export interface ChatMessage {
  id: string;
  role: 'user' | 'agent';
  blocks: ContentBlock[];
  timestamp: number;
}

export type ContentBlock =
  | { type: 'text'; content: string }
  | { type: 'thinking'; content: string }
  | { type: 'card'; card_type: CardType; data: unknown }
  | { type: 'tool_call'; name: string; status: 'running' | 'done' | 'error' };

// ===== API Request/Response =====

export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T | null;
}

export interface CreateSessionResponse {
  session_id: string;
  phase: Phase;
  created_at: string;
}

export interface ChatRequest {
  session_id: string;
  message: string;
}
