import { useCallback, useEffect, useRef } from 'react';
import { RiTelegram2Fill } from 'react-icons/ri';
import { Button } from 'ui/button';

import { TelegramInfo } from '../types';

interface TelegramWidgetProps {
    onAuth: (user: TelegramInfo) => void;
    disabled?: boolean;
    isLoading?: boolean;
    buttonText?: string;
}

export function TelegramWidget({ 
    onAuth, 
    disabled = false, 
    isLoading = false,
    buttonText = "Привязать Telegram"
}: TelegramWidgetProps) {
    const buttonRef = useRef<HTMLDivElement>(null);

    const onTelegramAuth = useCallback(
        (user: TelegramInfo) => {
            onAuth(user);
        },
        [onAuth],
    );

    useEffect(() => {
        if (disabled) return;

        // Создаем уникальное имя для функции, чтобы избежать конфликтов
        const callbackName = `onTelegramAuth_${Date.now()}`;
        
        // @ts-ignore
        window[callbackName] = onTelegramAuth;

        const script = document.createElement('script');
        script.src = 'https://telegram.org/js/telegram-widget.js?22';
        script.setAttribute('data-telegram-login', 'vkr_pulse_bot');
        script.setAttribute('data-size', 'large');
        script.setAttribute('data-userpic', 'false');
        script.setAttribute('data-radius', '0');
        script.setAttribute('data-onauth', `${callbackName}(user)`);
        script.async = true;

        if (buttonRef.current) {
            buttonRef.current.innerHTML = '';
            buttonRef.current.appendChild(script);
        }

        return () => {
            // @ts-ignore
            delete window[callbackName];
        };
    }, [onTelegramAuth, disabled]);

    return (
        <div className="relative">
            <Button
                type="button"
                disabled={disabled || isLoading}
                variant="default"
            >
                <RiTelegram2Fill />
                <span>{isLoading ? 'Привязываем...' : buttonText}</span>
            </Button>
            {!disabled && (
                <div
                    ref={buttonRef}
                    className="absolute top-0 left-0 w-full h-full z-10 opacity-0"
                />
            )}
        </div>
    );
} 