import useUpdateDocumentTitleOnce from './useUpdateDocumentTitleOnce.ts'
import {reduxRenderHook} from "../test-utils.tsx";
import {RootState} from "../state/store.ts";

describe('use Update Document Title Once', () => {
  it('should accept parameter that changes the document title', () => {
    const newTitle = 'new title'

    const [_, store] = reduxRenderHook(() => useUpdateDocumentTitleOnce(newTitle))

    expect((store.getState() as RootState).global.documentTitle).toBe(newTitle)
  })
})
