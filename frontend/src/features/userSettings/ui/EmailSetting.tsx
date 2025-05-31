import { zodResolver } from '@hookform/resolvers/zod';
import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { currentUserApi } from 'src/entities/user';
import { ApiErrorCode, isBaseError, useAppDispatch } from 'src/store';
import { Button } from 'ui/button';
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

import { userSettingApi, useSettingQuery, useUpdateEmailSettingMutation } from '../api';

const EmailSettingFormSchema = z.object({
    email: z.string().email('Введите корректный email'),
});

export function EmailSetting() {
    const dispatch = useAppDispatch();
    const { data } = useSettingQuery();
    const [updateEmail, { isLoading, error }] = useUpdateEmailSettingMutation();

    const form = useForm<z.infer<typeof EmailSettingFormSchema>>({
        resolver: zodResolver(EmailSettingFormSchema),
        defaultValues: {
            email: data?.email || '',
        },
    });

    const onUpdateEmail = (email: string) => {
        updateEmail({
            email,
        })
            .unwrap()
            .then(() => {
                dispatch(
                    userSettingApi.util.updateQueryData(
                        'setting',
                        undefined,
                        (draft) => {
                            draft.email = email;
                        },
                    )
                );
                dispatch(
                    currentUserApi.util.updateQueryData('getUser', undefined, (draft) => {
                        draft.email = email;
                    })
                );
                toast.success('Email успешно обновлен');
            })
            .catch(() => {
                toast.error('Произошла ошибка при обновлении email');
            });
    };

    const onSubmit = form.handleSubmit((data) => onUpdateEmail(data.email));

    useEffect(() => {
        if (!isBaseError(error)) return;

        if (error.data.error === ApiErrorCode.CodeAlreadyExists) {
            form.setError('email', {
                message: 'Почта уже зарегистрирована',
            });
        }
    }, [error, form]);

    if (!data) {
        return null;
    }

    return (
        <div className="flex flex-col gap-4 w-full">
            <div className="flex flex-col gap-4">
                <div className="flex flex-col">
                    <h3 className="text-lg font-medium">Email</h3>
                    <p className="text-sm text-muted-foreground">
                        Изменить адрес электронной почты для входа
                    </p>
                </div>
                
                <Form {...form}>
                    <form onSubmit={onSubmit} className="flex flex-col gap-4">
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
                        <div className="flex justify-end">
                            <Button type="submit" disabled={isLoading} className="w-fit">
                                Обновить email
                            </Button>
                        </div>
                    </form>
                </Form>
            </div>
        </div>
    );
}
