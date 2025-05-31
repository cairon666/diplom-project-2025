import { zodResolver } from '@hookform/resolvers/zod';
import { Link } from '@tanstack/react-router';
import { useForm } from 'react-hook-form';
import { router } from 'src/app/App';
import { Route as loginRouter } from 'src/app/routes/auth/login';
import { ApiErrorCode, isApiError, isFetchBaseQueryError } from 'src/store';
import { Button } from 'ui/button';
import {
    Card,
    CardHeader,
    CardTitle,
    CardDescription,
    CardContent,
    CardFooter,
} from 'ui/card';
import {
    FormField,
    FormItem,
    FormLabel,
    FormControl,
    FormMessage,
    Form,
} from 'ui/form';
import { Input } from 'ui/input';
import { Separator } from 'ui/separator';
import { z } from 'zod';

import { useRegisterMutation } from '../../api';
import { TelegramButton } from '../TelegramButton';

const RegisterFormSchema = z.object({
    first_name: z.string({
        required_error: 'Введите имя',
    }),
    second_name: z.string({
        required_error: 'Введите фамилию',
    }),
    email: z
        .string({
            required_error: 'Введите email',
        })
        .email('Введите корректный email'),
    password: z.string({
        required_error: 'Введите пароль',
    }),
});

export function Register() {
    const form = useForm<z.infer<typeof RegisterFormSchema>>({
        resolver: zodResolver(RegisterFormSchema),
        defaultValues: {
            first_name: '',
            second_name: '',
            email: '',
            password: '',
        },
    });

    const [register, { isLoading }] = useRegisterMutation();

    const onSubmit = form.handleSubmit((data) => {
        register({
            first_name: data.first_name,
            second_name: data.second_name,
            email: data.email,
            password: data.password,
        })
            .unwrap()
            .then(() =>
                router.navigate({
                    to: `${loginRouter.to}?email=${data.email}&password=${data.password}`,
                }),
            )
            .catch((error: unknown) => {
                if (isFetchBaseQueryError(error) && isApiError(error.data)) {
                    const errorCode = error.data.error;
                    if (errorCode === ApiErrorCode.CodeEmailAlreadyExists) {
                        form.setError('email', {
                            message:
                                'Пользователь с таким email уже существует',
                        });
                    }
                }
            });
    });

    return (
        <Form {...form}>
            <form onSubmit={onSubmit}>
                <Card className="w-100 space-1.5">
                    <CardHeader>
                        <CardTitle>Регистрация</CardTitle>
                        <CardDescription>
                            Пожалуйста, введите данные для регистрации
                        </CardDescription>
                    </CardHeader>
                    <CardContent>
                        <TelegramButton />
                        <div className="relative">
                            <Separator className="mb-4" />
                            <span className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 bg-white px-2">
                                или
                            </span>
                        </div>
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
                                    name="second_name"
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
                            type="submit"
                            disabled={isLoading}
                            className="w-full mb-2"
                        >
                            Зарегистрироваться
                        </Button>
                        <span className="text-center">
                            <span>У вас есть аккаунт? </span>
                            <Link
                                className="underline text-blue-600 hover:text-blue-800 visited:text-purple-600"
                                to={loginRouter.to}
                            >
                                Войти
                            </Link>
                        </span>
                    </CardFooter>
                </Card>
            </form>
        </Form>
    );
}
