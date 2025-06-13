import type { RootState } from 'src/store/store';

import { Store, configureStore } from '@reduxjs/toolkit';
import { render } from '@testing-library/react';
import { Provider } from 'react-redux';
import { DeepPartial } from 'utility-types';

// Создаем минимальный reducer для тестов
const createTestReducer =
    () =>
    (state = {}) =>
        state;

// Функция для создания тестового стора
export const createTestStore = (preloadedState?: DeepPartial<RootState>) => {
    return configureStore({
        reducer: createTestReducer(),
        preloadedState: preloadedState as any,
        middleware: (getDefaultMiddleware) =>
            getDefaultMiddleware({
                serializableCheck: false,
                immutableCheck: false,
            }),
        // Отключаем DevTools для тестов
        devTools: false,
    });
};

// Функция для рендеринга с готовым стором
export const renderWithStore = (
    component: React.ReactElement,
    store: Store,
) => {
    return render(<Provider store={store}>{component}</Provider>);
};

// Функция для рендеринга с автоматическим созданием стора
export const renderWithTestStore = (
    component: React.ReactElement,
    preloadedState?: DeepPartial<RootState>,
) => {
    const store = createTestStore(preloadedState);
    return {
        ...render(<Provider store={store}>{component}</Provider>),
        store,
    };
};

// Функция для рендеринга с начальным состоянием (сокращенная запись)
export const renderWithInitialState = (
    component: React.ReactElement,
    initialState?: DeepPartial<RootState>,
) => {
    return renderWithTestStore(component, initialState);
};
