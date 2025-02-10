import { reduxRender } from '../test-utils.tsx'
import NotFound from './NotFound.tsx'
import { screen } from '@testing-library/react'
import { RootState } from '../state/store.ts'

describe('Not Found', () => {
  it('should render', () => {
    const [_, store] = reduxRender(<NotFound />)

    expect((store.getState() as RootState).global.documentTitle).toMatch(/not found/i)
    expect(screen.getByRole('heading', { name: /not found/i })).toBeInTheDocument()
  })
})
