import styled from '@emotion/styled';
import type { FC } from 'react';

interface MessageProps {
  author: string;
  content: string;
}

const MessageContainer = styled.div`
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

export const Message: FC<MessageProps> = (props: MessageProps) => {
  const { author, content } = props;

  return (
    <MessageContainer>
      <Author>{author}</Author>
      <Content>{content}</Content>
    </MessageContainer>
  );
};
