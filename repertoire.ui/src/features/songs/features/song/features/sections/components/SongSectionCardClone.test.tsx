import {
  emptyArtist,
  emptySong,
  emptySongPart,
  emptySongSection,
  mantineRender
} from '../../../../../../../test-utils.tsx'
import SongSectionCardClone from './SongSectionCardClone.tsx'
import { Instrument, SongSection } from '../../../../../../../types/models/Song.ts'
import { BandMember } from '../../../../../../../types/models/Artist.ts'
import { screen, within } from '@testing-library/react'
import { SongProvider } from '../../../context/SongContext.tsx'

describe('Song Section Card Clone', () => {
  const section: SongSection = {
    ...emptySongSection,
    id: '1',
    name: 'Solo 1',
    rehearsals: 12,
    confidence: 50,
    progress: 150,
    songSectionType: {
      id: 'some-id',
      name: 'Solo'
    }
  }

  function render(
    isDragging = false,
    isDropAnimating = false,
    song = emptySong,
    cloneSection: SongSection = section
  ) {
    return mantineRender(
      <SongProvider song={song}>
        <SongSectionCardClone
          section={cloneSection}
          isDragging={isDragging}
          isDropAnimating={isDropAnimating}
          maxSectionProgress={200}
        />
      </SongProvider>
    )
  }

  it('should render and display section information', () => {
    render()

    expect(screen.getByLabelText(`song-section-${section.name}-clone`)).toBeInTheDocument()
    expect(screen.getByText(section.name)).toBeInTheDocument()
    expect(screen.getByText(section.songSectionType.name)).toBeInTheDocument()
    expect(screen.getByText(section.rehearsals)).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'confidence' })).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'confidence' })).toHaveValue(section.confidence)
    expect(screen.getByRole('progressbar', { name: 'progress' })).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'progress' })).toHaveValue(75)
  })

  it('should render aggregated instruments and band members', () => {
    const bandMembers: BandMember[] = [
      {
        id: '1',
        name: 'Mike',
        roles: [{ id: '1', name: 'Guitarist' }],
        imageUrl: 'default.png'
      },
      {
        id: '2',
        name: 'Leonard',
        roles: [{ id: '2', name: 'Voice' }]
      }
    ]
    const instruments: Instrument[] = [
      { id: '1', name: 'Electric Guitar' },
      { id: '2', name: 'Voice' }
    ]

    render(
      false,
      false,
      { ...emptySong, artist: { ...emptyArtist, isBand: true } },
      {
        ...section,
        parts: [
          {
            ...emptySongPart,
            id: '1',
            instrument: instruments[0],
            bandMembers: bandMembers
          },
          {
            ...emptySongPart,
            id: '2',
            instrument: instruments[1],
            bandMembers: [bandMembers[1]]
          }
        ]
      }
    )

    const instrumentsEl = screen.getByLabelText('instruments')
    const bandMembersEl = screen.getByLabelText('band-members')

    instruments.forEach((instrument) => {
      expect(within(instrumentsEl).getByLabelText(instrument.name)).toBeInTheDocument()
    })
    bandMembers.forEach((bandMember) => {
      if (bandMember.imageUrl) {
        expect(
          within(bandMembersEl).getByRole('img', { name: bandMember.name })
        ).toBeInTheDocument()
      } else {
        expect(
          within(bandMembersEl).getByLabelText(`default-icon-${bandMember.name}`)
        ).toBeInTheDocument()
      }
    })
  })

  it('should not render band members for an artist that is not a band', () => {
    render(false, false, { ...emptySong, artist: { ...emptyArtist, isBand: false } })

    expect(screen.queryByLabelText('band-members')).not.toBeInTheDocument()
  })

  it('should be selected while dragging', async () => {
    render(true)

    expect(screen.getByLabelText(`song-section-${section.name}-clone`)).toHaveAttribute(
      'aria-selected',
      'true'
    )
  })
})
