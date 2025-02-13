import { reduxRender } from '../test-utils.tsx'
import Home from './Home.tsx'
import { screen } from '@testing-library/react'
import { expect } from 'vitest'
import { RootState } from '../state/store.ts'

describe('Home', () => {
  it('should render', () => {
    const [_, store] = reduxRender(<Home />)

    expect((store.getState() as RootState).global.documentTitle).toMatch(/home/i)
    expect(screen.getByLabelText('top')).toBeInTheDocument()
    expect(screen.getByLabelText('genres')).toBeInTheDocument()
    expect(screen.getByLabelText('recently-played-songs')).toBeInTheDocument()
    expect(screen.getByLabelText('playlists')).toBeInTheDocument()
    expect(screen.getByLabelText('top-artists')).toBeInTheDocument()
  })
})
