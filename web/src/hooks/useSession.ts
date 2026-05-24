import { useCallback, useEffect } from 'react';
import { createSession, deleteSession } from '../api/client';
import { useChatStore } from '../stores/chatStore';

export function useSession() {
  const { sessionId, setSessionId, setPhase, reset } = useChatStore();

  const initSession = useCallback(async () => {
    try {
      const data = await createSession();
      setSessionId(data.session_id);
      setPhase(data.phase);
    } catch (err) {
      console.error('Failed to create session:', err);
    }
  }, [setSessionId, setPhase]);

  const endSession = useCallback(async () => {
    if (!sessionId) return;
    try {
      await deleteSession(sessionId);
    } catch {
      // Ignore errors on cleanup
    } finally {
      reset();
    }
  }, [sessionId, reset]);

  // Create session on mount
  useEffect(() => {
    if (!sessionId) {
      initSession();
    }
  }, [sessionId, initSession]);

  // Clean up on unmount
  useEffect(() => {
    return () => {
      // Don't delete session on unmount — user may refresh and want to continue
    };
  }, []);

  return { sessionId, initSession, endSession };
}
