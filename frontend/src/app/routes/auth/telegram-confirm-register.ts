import { createFileRoute } from '@tanstack/react-router';
import { TelegramConfirmRegister } from 'src/features/auth';

export const Route = createFileRoute('/auth/telegram-confirm-register')({
    component: TelegramConfirmRegister,
});
