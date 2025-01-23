import { reduxRender } from '../test-utils.tsx'
import Unauthorized from './Unauthorized.tsx'
import { screen } from '@testing-library/react'
import { RootState } from '../state/store.ts'

describe('Unauthorized', () => {
  it('should render', () => {
    const [_, store] = reduxRender(<Unauthorized />)

    expect((store.getState() as RootState).global.documentTitle).toMatch(/unauthorized/i)
    expect(screen.getByRole('heading', { name: /unauthorized/i })).toBeInTheDocument()
  })
})
