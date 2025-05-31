import { zodResolver } from '@hookform/resolvers/zod';
import { useSearch, useNavigate, Navigate } from '@tanstack/react-router';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { Route as loginRouter } from 'src/app/routes/auth/login';
import { ApiErrorCode, isBaseError } from 'src/store';
import { Button } from 'ui/button';
import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from 'ui/card';
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from 'ui/form';
import { Input } from 'ui/input';
import { z } from 'zod';

import { useTelegramConfirmMutation } from '../../api';
import { useOnSuccessLogin } from '../../hooks/useOnSuccessLogin';

const TelegramConfirmFormSchema = z.object({
    email: z
        .string({
            required_error: 'Введите email',
        })
        .email('Введите корректный email'),
    first_name: z.string({
        required_error: 'Введите имя',
    }),
    last_name: z.string({
        required_error: 'Введите фамилию',
    }),
});

interface TelegramConfirmRegisterSearchParams {
    tempId: string;
    first_name: string;
    last_name: string;
}

export function TelegramConfirmRegister() {
    const navigate = useNavigate();
    const search: TelegramConfirmRegisterSearchParams = useSearch({
        strict: false,
    });

    const form = useForm<z.infer<typeof TelegramConfirmFormSchema>>({
        resolver: zodResolver(TelegramConfirmFormSchema),
        defaultValues: {
            email: '',
            first_name: search.first_name || '',
            last_name: search.last_name || '',
        },
    });

    const [telegramConformMutation, { isLoading }] =
        useTelegramConfirmMutation();
    const { onSuccessLogin } = useOnSuccessLogin();

    const onSubmit = form.handleSubmit((data) => {
        telegramConformMutation({
            first_name: data.first_name,
            last_name: data.last_name,
            email: data.email,
            temp_id: search.tempId,
        })
            .unwrap()
            .then(onSuccessLogin)
            .catch((error: unknown) => {
                if (isBaseError(error)) {
                    const errorCode = error.data.error;
                    if (errorCode === ApiErrorCode.CodeEmailAlreadyExists) {
                        form.setError('email', {
                            message:
                                'Пользователь с таким email уже существует',
                        });
                    }

                    if (
                        errorCode === ApiErrorCode.CodeTempIdNotFound ||
                        errorCode === ApiErrorCode.CodeInternalError
                    ) {
                        toast.error('Что-то пошло не так. Попробуйте позже');
                        navigate({ to: loginRouter.to });
                    }
                }
            });
    });

    if (!search.tempId) {
        return <Navigate to={loginRouter.to} />;
    }

    return (
        <Form {...form}>
            <form onSubmit={onSubmit}>
                <Card className="w-100 space-1.5">
                    <CardHeader>
                        <CardTitle>Подтвердите регистрацию</CardTitle>
                        <CardDescription>
                            Пожалуйста, введите свои данные
                        </CardDescription>
                    </CardHeader>
                    <CardContent>
                        <div className="grid w-full items-center gap-4">
                            <div className="flex gap-2">
                                <FormField
                                    control={form.control}
                                    name="first_name"
                                    render={({ field }) => (
                                        <FormItem className="flex flex-col gap-1">
                                            <FormLabel>Имя</FormLabel>
                                            <FormControl>
                                                <Input
                                                    {...field}
                                                    placeholder="Введите имя"
                                                />
                                            </FormControl>
                                            <FormMessage />
                                        </FormItem>
                                    )}
                                />
                                <FormField
                                    control={form.control}
                                    name="last_name"
                                    render={({ field }) => (
                                        <FormItem className="flex flex-col gap-1">
                                            <FormLabel>Фамилия</FormLabel>
                                            <FormControl>
                                                <Input
                                                    {...field}
                                                    placeholder="Введите фамилию"
                                                />
                                            </FormControl>
                                            <FormMessage />
                                        </FormItem>
                                    )}
                                />
                            </div>
                            <FormField
                                control={form.control}
                                name="email"
                                render={({ field }) => (
                                    <FormItem className="flex flex-col gap-1">
                                        <FormLabel>Почта</FormLabel>
                                        <FormControl>
                                            <Input
                                                {...field}
                                                placeholder="Введите почту"
                                            />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                        </div>
                    </CardContent>
                    <CardFooter className="flex-col">
                        <Button
                            disabled={isLoading}
                            type="submit"
                            className="w-full mb-2"
                        >
                            Подтвердить
                        </Button>
                    </CardFooter>
                </Card>
            </form>
        </Form>
    );
}
