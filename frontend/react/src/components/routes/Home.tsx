import type { FC, FormEvent } from 'react';
import { useAuth } from 'react-oidc-context';
import { useEffect, useState } from 'react';
import { Message } from '../Message.tsx';

interface MessageData {
  id: string;
  author: string;
  content: string;
  createdAt: string;
}

export const Home: FC = () => {
  const auth = useAuth();
  const [messages, setMessages] = useState<MessageData[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [newMessage, setNewMessage] = useState<string>('');
  const [isSending, setIsSending] = useState<boolean>(false);

  const fetchMessages = async () => {
    try {
      setLoading(true);
      setError(null);

      const response = await fetch('http://localhost:8080/messages?limit=100&offset=0', {
        headers: {
          Authorization: `Bearer ${auth.user?.access_token}`,
        },
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch messages: ${response.statusText}`);
      }

      const data = await response.json();
      setMessages(data);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchMessages();
  }, [auth.user?.access_token]);

  const handleDeleteMessage = async (id: string) => {
    try {
      const response = await fetch(`http://localhost:8080/messages/${id}`, {
        method: 'DELETE',
        headers: {
          Authorization: `Bearer ${auth.user?.access_token}`,
        },
      });

      if (!response.ok) {
        throw new Error(`Failed to delete message: ${response.statusText}`);
      }

      // Remove the deleted message from the state
      setMessages((prevMessages) => prevMessages.filter((message) => message.id !== id));
    } catch (err: any) {
      setError(err.message);
    }
  };

  const handleSendMessage = async (e: FormEvent) => {
    e.preventDefault();

    try {
      setIsSending(true);
      const author = auth.user?.profile.preferred_username || auth.user?.profile.email || 'Unknown';

      const optimisticMessage = {
        id: `temp-${Date.now()}`,
        author,
        content: newMessage,
        createdAt: new Date().toISOString(),
      };
      setMessages((prevMessages) => [optimisticMessage, ...prevMessages]);

      const response = await fetch('http://localhost:8080/messages', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${auth.user?.access_token}`,
        },
        body: JSON.stringify({ author, content: newMessage }),
      });

      if (!response.ok) {
        throw new Error(`Failed to send message: ${response.statusText}`);
      }

      const createdMessage = await response.json();
      setMessages((prevMessages) =>
        prevMessages.map((message) =>
          message.id === optimisticMessage.id ? createdMessage : message
        )
      );

      setNewMessage('');
      fetchMessages();
    } catch (err: any) {
      setError(err.message);
    } finally {
      setIsSending(false);
    }
  };

  if (loading) {
    return <p>Loading messages...</p>;
  }

  if (error) {
    return <p>Error: {error}</p>;
  }

  const loggedInUser = auth.user?.profile.preferred_username || auth.user?.profile.email;

  return (
    <>
      <h1>Home</h1>
      <p>Write a few messages...</p>

      <h2>Messages</h2>
      {messages.length > 0 ? (
        messages.map((message) => (
          <Message
            key={message.id}
            author={message.author}
            content={message.content}
            isAuthor={message.author === loggedInUser} // Check if the logged-in user is the author
            onDelete={() => handleDeleteMessage(message.id)} // Pass the delete handler
          />
        ))
      ) : (
        <p>No messages found.</p>
      )}

      <h2>Send a Message</h2>
      <form onSubmit={handleSendMessage}>
        <textarea
          value={newMessage}
          onChange={(e) => setNewMessage(e.target.value)}
          placeholder="Write your message here..."
          rows={4}
          style={{ width: '100%', marginBottom: '8px' }}
          disabled={isSending}
        />
        <button type="submit" disabled={!newMessage.trim() || isSending}>
          {isSending ? 'Sending...' : 'Send'}
        </button>
      </form>
    </>
  );
};