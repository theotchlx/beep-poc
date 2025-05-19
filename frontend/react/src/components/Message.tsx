import styled from '@emotion/styled';
import type { FC } from 'react';

interface MessageProps {
  author: string;
  content: string;
  isAuthor: boolean; // New prop to determine if the logged-in user is the author
  onDelete?: () => void; // Callback for deleting the message
}

const MessageContainer = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  border: 1px solid #ccc;
  border-radius: 8px;
  padding: 16px;
  margin: 8px 0;
  background-color: #f9f9f9;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
`;

const Author = styled.div`
  font-weight: bold;
  margin-bottom: 8px;
  color: #333;
`;

const Content = styled.div`
  font-size: 14px;
  color: #555;
`;

const DeleteButton = styled.button`
  background-color: #ff4d4d;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 4px 8px;
  cursor: pointer;

  &:hover {
    background-color: #ff1a1a;
  }
`;

export const Message: FC<MessageProps> = ({ author, content, isAuthor, onDelete }) => {
  return (
    <MessageContainer>
      <div>
        <Author>{author}</Author>
        <Content>{content}</Content>
      </div>
      {isAuthor && (
        <DeleteButton onClick={onDelete}>
          Delete
        </DeleteButton>
      )}
    </MessageContainer>
  );
};