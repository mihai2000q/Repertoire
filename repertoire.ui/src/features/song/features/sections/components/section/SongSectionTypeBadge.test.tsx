import { screen } from '@testing-library/react'
import { mantineRender } from '../../../../../../test-utils.tsx'
import SongSectionTypeBadge from './SongSectionTypeBadge.tsx'
import { SongSectionType } from '../../../../../../types/models/Song.ts'

describe('Rehearsals Badge', () => {
  it('should render', () => {
    const songSectionType: SongSectionType = {
      id: '1',
      name: 'Verse'
    }

    mantineRender(<SongSectionTypeBadge songSectionType={songSectionType} />)

    expect(screen.getByText(songSectionType.name)).toBeInTheDocument()
  })
})
