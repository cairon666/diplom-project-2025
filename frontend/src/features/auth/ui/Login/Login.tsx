import { zodResolver } from '@hookform/resolvers/zod';
import { Link, useSearch } from '@tanstack/react-router';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { Route as registerRouter } from 'src/app/routes/auth/register';
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
import { Separator } from 'ui/separator';
import { z } from 'zod';

import { useLoginMutation } from '../../api';
import { useOnSuccessLogin } from '../../hooks/useOnSuccessLogin';
import { TelegramButton } from '../TelegramButton';

const LoginFormSchema = z.object({
    email: z
        .string({
            required_error: 'Введите email',
        })
        .email('Введите корректный email'),
    password: z.string({
        required_error: 'Введите пароль',
    }),
});

interface LoginSearchParams {
    email?: string;
    password?: string;
    redirect?: string;
}

export function Login() {
    const search: LoginSearchParams = useSearch({ strict: false });
    const form = useForm<z.infer<typeof LoginFormSchema>>({
        resolver: zodResolver(LoginFormSchema),
        defaultValues: {
            email: search.email || '',
            password: search.password || '',
        },
    });
    const { onSuccessLogin } = useOnSuccessLogin();
    const [login, { isLoading }] = useLoginMutation();

    const onSubmitError = (error: unknown) => {
        if (isBaseError(error)) {
            const errorCode = error.data.error;
            if (errorCode === ApiErrorCode.CodeLoginNotRegistered) {
                form.setError('email', {
                    message: 'Почта не зарегистрирована',
                });
            } else if (errorCode === ApiErrorCode.CodeWrongPassword) {
                form.setError('password', {
                    message: 'Неверный пароль',
                });
            } else {
                toast.error('Что-то пошло не так. Попробуйте позже');
            }
        }
    };

    const onSubmit = form.handleSubmit((data) => {
        login({
            email: data.email,
            password: data.password,
        })
            .unwrap()
            .then(onSuccessLogin)
            .catch(onSubmitError);
    });

    return (
        <Form {...form}>
            <form onSubmit={onSubmit}>
                <Card className="w-100 space-1.5">
                    <CardHeader>
                        <CardTitle>Вход</CardTitle>
                        <CardDescription>
                            Пожалуйста, войдите в свой аккаунт
                        </CardDescription>
                    </CardHeader>
                    <CardContent>
                        <TelegramButton />
                        <div className="relative">
                            <Separator className="mb-2" />
                            <span className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 bg-white px-2">
                                или
                            </span>
                        </div>
                        <div className="grid w-full items-center gap-4">
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
                            <FormField
                                control={form.control}
                                name="password"
                                render={({ field }) => (
                                    <FormItem className="flex flex-col gap-1">
                                        <FormLabel>Пароль</FormLabel>
                                        <FormControl>
                                            <Input
                                                {...field}
                                                placeholder="Введите пароль"
                                                type="password"
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
                            Войти
                        </Button>
                        <span className="text-center">
                            <span>У вас нет учетной записи? </span>
                            <Link
                                className="underline text-blue-600 hover:text-blue-800 visited:text-purple-600"
                                to={registerRouter.to}
                            >
                                Зарегистрируйтесь
                            </Link>
                        </span>
                    </CardFooter>
                </Card>
            </form>
        </Form>
    );
}
