import { UUID } from '@/types/GenericTypes';
import { CardDto } from '@/models/CardDto';
import { RoundState } from '@/types/RoundState';

interface CommonMessage {
  id: UUID;
  type: string;
  payload: object;
}

export interface RefreshMessage extends CommonMessage {
  type: 'refresh';
  payload: {
    card: CardDto;
  };
}

export interface RoundMessage extends CommonMessage {
  type: 'round';
  payload: {
    isWin: boolean;
    card: CardDto;
    round_state?: RoundState;
  };
}

export interface CounterMessage extends CommonMessage {
  type: 'counter';
  payload: {
    counter: number;
  };
}

export interface UpdateGameStateMessage extends CommonMessage {
  type: 'state';
  payload: {
    counter?: number;
    cards: CardDto[];
    round_state?: RoundState;
  };
}

export type CommunicationMessages = RoundMessage | RefreshMessage | CounterMessage | UpdateGameStateMessage;
