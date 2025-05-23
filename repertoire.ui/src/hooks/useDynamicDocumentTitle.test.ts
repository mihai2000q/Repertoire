import { act } from '@testing-library/react'
import useDynamicDocumentTitle from './useDynamicDocumentTitle.ts'
import {reduxRenderHook} from "../test-utils.tsx";
import {RootState} from "../state/store.ts";

describe('use Dynamic Document Title', () => {
  it('should return a setter that changes the document title', () => {
    const newTitle = 'new title'
    const secondNewTitle = 'second'
    const thirdNewTitle = 'third'

    const [{ result }, store] = reduxRenderHook(() => useDynamicDocumentTitle())

    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    act(() => result.current(newTitle))
    expect((store.getState() as RootState).global.documentTitle).toBe(newTitle)

    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    act(() => result.current(secondNewTitle))
    expect((store.getState() as RootState).global.documentTitle).toBe(secondNewTitle)

    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    act(() => result.current((prevTitle: string) => prevTitle + thirdNewTitle))
    expect((store.getState() as RootState).global.documentTitle).toBe(secondNewTitle + thirdNewTitle)
  })
})
