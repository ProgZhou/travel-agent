import type { ApiResponse, CreateSessionResponse, Session } from '../types';

const BASE_URL = '/api';

class ApiError extends Error {
  constructor(
    public code: number,
    message: string,
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

async function request<T>(
  path: string,
  options?: RequestInit,
): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    headers: { 'Content-Type': 'application/json', ...options?.headers },
    ...options,
  });

  const json: ApiResponse<T> = await res.json();

  if (json.code !== 0) {
    throw new ApiError(json.code, json.message);
  }

  return json.data as T;
}

export async function createSession(): Promise<CreateSessionResponse> {
  return request<CreateSessionResponse>('/session', {
    method: 'POST',
    body: '{}',
  });
}

export async function getSession(id: string): Promise<Session> {
  return request<Session>(`/session/${id}`);
}

export async function deleteSession(id: string): Promise<void> {
  await request<null>(`/session/${id}`, { method: 'DELETE' });
}

export { ApiError };
