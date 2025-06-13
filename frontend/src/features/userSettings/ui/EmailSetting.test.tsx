import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { toast } from 'sonner';
import { createTestStore, renderWithStore } from 'src/shared/lib/test';
import { vi, describe, it, expect, beforeEach } from 'vitest';

import { useSettingQuery, useUpdateEmailSettingMutation } from '../api';
import { EmailSetting } from './EmailSetting';

// Моки должны быть в начале файла без использования переменных
vi.mock('sonner', () => ({
    toast: {
        success: vi.fn(),
        error: vi.fn(),
    },
}));

vi.mock('src/store', () => ({
    useAppDispatch: () => vi.fn(),
    isBaseError: vi.fn(),
    ApiErrorCode: {
        CodeAlreadyExists: 'CODE_ALREADY_EXISTS',
    },
}));

vi.mock('../api', () => ({
    userSettingApi: {
        util: {
            updateQueryData: vi.fn(),
        },
    },
    useSettingQuery: vi.fn(),
    useUpdateEmailSettingMutation: vi.fn(),
}));

vi.mock('src/entities/user', () => ({
    currentUserApi: {
        util: {
            updateQueryData: vi.fn(),
        },
    },
}));

describe('EmailSetting', () => {
    const mockUpdateEmailMutation = vi.fn();
    let store: ReturnType<typeof createTestStore>;

    beforeEach(() => {
        vi.clearAllMocks();

        // Создаем новый тестовый стор для каждого теста
        store = createTestStore();

        // Настраиваем моки для каждого теста с полными типами
        vi.mocked(useSettingQuery).mockReturnValue({
            data: { email: 'test@example.com' },
            isLoading: false,
            isError: false,
            error: undefined,
            refetch: vi.fn(),
        } as any);

        vi.mocked(useUpdateEmailSettingMutation).mockReturnValue([
            mockUpdateEmailMutation,
            {
                isLoading: false,
                error: null,
                reset: vi.fn(),
            } as any,
        ]);
    });

    it('should render email setting form', () => {
        renderWithStore(<EmailSetting />, store);

        expect(screen.getByText('Email')).toBeInTheDocument();
        expect(
            screen.getByText('Изменить адрес электронной почты для входа'),
        ).toBeInTheDocument();
        expect(screen.getByLabelText('Почта')).toBeInTheDocument();
        expect(
            screen.getByRole('button', { name: 'Обновить email' }),
        ).toBeInTheDocument();
    });

    it('should populate form with current email', () => {
        renderWithStore(<EmailSetting />, store);

        const emailInput = screen.getByDisplayValue('test@example.com');
        expect(emailInput).toBeInTheDocument();
    });

    it('should not render if no data available', () => {
        vi.mocked(useSettingQuery).mockReturnValue({
            data: null,
            isLoading: false,
            isError: false,
            error: undefined,
            refetch: vi.fn(),
        } as any);

        const { container } = renderWithStore(<EmailSetting />, store);
        expect(container.firstChild).toBeNull();
    });

    it('should show validation error for invalid email', async () => {
        const user = userEvent.setup();
        renderWithStore(<EmailSetting />, store);

        const emailInput = screen.getByLabelText('Почта');
        const submitButton = screen.getByRole('button', {
            name: 'Обновить email',
        });

        // Очищаем поле и вводим невалидный email
        await user.clear(emailInput);
        await user.type(emailInput, 'invalid-email');
        await user.click(submitButton);

        await waitFor(() => {
            expect(
                screen.getByText('Введите корректный email'),
            ).toBeInTheDocument();
        });
    });

    it('should submit form with valid email', async () => {
        const user = userEvent.setup();
        const unwrapMock = vi.fn().mockResolvedValue({});

        mockUpdateEmailMutation.mockReturnValue({
            unwrap: unwrapMock,
        });

        renderWithStore(<EmailSetting />, store);

        const emailInput = screen.getByLabelText('Почта');
        const submitButton = screen.getByRole('button', {
            name: 'Обновить email',
        });

        // Очищаем поле и вводим новый email
        await user.clear(emailInput);
        await user.type(emailInput, 'newemail@example.com');
        await user.click(submitButton);

        await waitFor(() => {
            expect(mockUpdateEmailMutation).toHaveBeenCalledWith({
                email: 'newemail@example.com',
            });
        });
    });

    it('should show success toast on successful email update', async () => {
        const user = userEvent.setup();
        const unwrapMock = vi.fn().mockResolvedValue({});

        mockUpdateEmailMutation.mockReturnValue({
            unwrap: unwrapMock,
        });

        renderWithStore(<EmailSetting />, store);

        const emailInput = screen.getByLabelText('Почта');
        const submitButton = screen.getByRole('button', {
            name: 'Обновить email',
        });

        await user.clear(emailInput);
        await user.type(emailInput, 'newemail@example.com');
        await user.click(submitButton);

        await waitFor(() => {
            expect(toast.success).toHaveBeenCalledWith(
                'Email успешно обновлен',
            );
        });
    });

    it('should show error toast on failed email update', async () => {
        const user = userEvent.setup();
        const unwrapMock = vi.fn().mockRejectedValue(new Error('API Error'));

        mockUpdateEmailMutation.mockReturnValue({
            unwrap: unwrapMock,
        });

        renderWithStore(<EmailSetting />, store);

        const emailInput = screen.getByLabelText('Почта');
        const submitButton = screen.getByRole('button', {
            name: 'Обновить email',
        });

        await user.clear(emailInput);
        await user.type(emailInput, 'newemail@example.com');
        await user.click(submitButton);

        await waitFor(() => {
            expect(toast.error).toHaveBeenCalledWith(
                'Произошла ошибка при обновлении email',
            );
        });
    });

    it('should disable submit button when loading', () => {
        vi.mocked(useUpdateEmailSettingMutation).mockReturnValue([
            mockUpdateEmailMutation,
            {
                isLoading: true,
                error: null,
                reset: vi.fn(),
            } as any,
        ]);

        renderWithStore(<EmailSetting />, store);

        const submitButton = screen.getByRole('button', {
            name: 'Обновить email',
        });
        expect(submitButton).toBeDisabled();
    });
});
