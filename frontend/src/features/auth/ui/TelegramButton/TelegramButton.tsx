import { useNavigate } from '@tanstack/react-router';
import { useCallback, useEffect, useRef } from 'react';
import { RiTelegram2Fill } from 'react-icons/ri';
import { toast } from 'sonner';
import { Route as telegramConfigrmRegisterRouter } from 'src/app/routes/auth/telegram-confirm-register';
import {
    ApiErrorCode,
    AppErrorWithFields,
    isApiErrorWithFields,
    isBaseError,
} from 'src/store';
import { Button } from 'ui/button';

import { useTelegramLoginMutation } from '../../api';
import { useOnSuccessLogin } from '../../hooks/useOnSuccessLogin';

interface TelegramInfo {
    auth_date: number;
    first_name: string;
    last_name: string;
    username: string;
    hash: string;
    photo_url: string;
    id: number;
}

interface NeedRegistrationErrorProps {
    tempId: string;
}

export function TelegramButton() {
    const navigate = useNavigate();
    const buttonRef = useRef<HTMLDivElement>(null);

    const [telegramLoginMutation] = useTelegramLoginMutation();
    const { onSuccessLogin } = useOnSuccessLogin();

    const onNeedRegistrationError = useCallback(
        (
            error: AppErrorWithFields<NeedRegistrationErrorProps>,
            firstName: string,
            lastName: string,
        ) => {
            const { tempId } = error.fields;

            navigate({
                to: `${telegramConfigrmRegisterRouter.to}?tempId=${tempId}&first_name=${firstName}&last_name=${lastName}`,
            });
        },
        [navigate],
    );

    const onTelegramAuth = useCallback(
        (user: TelegramInfo) => {
            telegramLoginMutation(user)
                .unwrap()
                .then(onSuccessLogin)
                .catch((error: unknown) => {
                    if (isApiErrorWithFields(error)) {
                        if (
                            error.data.error ===
                            ApiErrorCode.CodeNeedEndRegistration
                        ) {
                            onNeedRegistrationError(
                                error.data as AppErrorWithFields<NeedRegistrationErrorProps>,
                                user.first_name,
                                user.last_name,
                            );
                        }
                    }

                    if (isBaseError(error)) {
                        if (
                            error.data.error ===
                            ApiErrorCode.CodeInvalidTelegramHash
                        ) {
                            toast.error(
                                'Хэш Telegram недействителен, пожалуйста, попробуйте снова',
                            );
                        }
                    }
                });
        },
        [onNeedRegistrationError, onSuccessLogin, telegramLoginMutation],
    );

    useEffect(() => {
        // @ts-ignore
        window.onTelegramAuth = onTelegramAuth;

        const script = document.createElement('script');
        script.src = 'https://telegram.org/js/telegram-widget.js?22';
        script.setAttribute('data-telegram-login', 'vkr_pulse_bot'); // Замените на имя вашего бота
        script.setAttribute('data-size', 'large');
        script.setAttribute('data-userpic', 'false');
        script.setAttribute('data-radius', '0');
        script.setAttribute('data-onauth', 'onTelegramAuth(user)');
        script.async = true;

        if (buttonRef.current) {
            buttonRef.current.innerHTML = '';
            buttonRef.current.appendChild(script);
        }

        return () => {
            // @ts-ignore
            delete window.onTelegramAuth;
        };
    }, [onTelegramAuth]);

    return (
        <div className="flex justify-space-between mb-4 gap-2 relative">
            <Button
                id="telegram-button"
                type="button"
                className="flex-1"
                variant="outline"
            >
                <RiTelegram2Fill />
                <span>Telegram</span>
            </Button>
            <div
                ref={buttonRef}
                className="absolute top-0 left-0 w-full h-full z-10 opacity-0"
            />
        </div>
    );
}
