import {
  emptyArtist,
  emptySong,
  emptySongSection,
  emptySongSettings,
  reduxRender,
  withToastify
} from '../../../../../../../../test-utils.tsx'
import { SongProvider } from '../../../../context/SongContext.tsx'
import { ReactNode } from 'react'
import AddNewSongSectionPart from './AddNewSongSectionPart.tsx'
import { screen } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { userEvent } from '@testing-library/user-event'
import Song, { Instrument, SongSection } from '../../../../../../../../types/models/Song.ts'
import { CreateSongPartRequest } from '../../../parts/types/requests/SongPartRequests.ts'
import { BandMember } from '../../../../../../../../types/models/Artist.ts'

describe('Add New Song Section Part', () => {
  const instruments: Instrument[] = [
    { id: '1', name: 'Guitar' },
    { id: '2', name: 'Piano' }
  ]

  const bandMembers: BandMember[] = [
    { id: '1', name: 'Chester', roles: [] },
    { id: '2', name: 'Michael', roles: [] }
  ]

  const section: SongSection = {
    ...emptySongSection,
    id: 'section-1',
    name: 'Chorus'
  }

  const handlers = [http.get('/songs/instruments', () => HttpResponse.json(instruments))]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  function render(ui: ReactNode, song: Song = emptySong) {
    return reduxRender(<SongProvider song={song}>{ui}</SongProvider>)
  }

  it('should render the add card', () => {
    render(<AddNewSongSectionPart section={section} />)

    expect(screen.getByLabelText(`add-new-song-section-part-card-${section.name}`)).toBeInTheDocument()
  })

  it('should open the form and focus the name input when the add card is clicked', async () => {
    const user = userEvent.setup()

    render(<AddNewSongSectionPart section={section} />)

    await user.click(screen.getByLabelText(`add-new-song-section-part-card-${section.name}`))

    expect(screen.getByRole('button', { name: 'close' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'select-band-member' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'select-band-member' })).not.toBeDisabled()
    expect(screen.getByRole('button', { name: 'select-instrument' })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /name/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /name/i })).toHaveFocus()
    expect(screen.getByRole('button', { name: 'Add' })).toBeInTheDocument()
  })

  it('should close the form when the close button is clicked', async () => {
    const user = userEvent.setup()

    render(<AddNewSongSectionPart section={section} />)

    await user.click(screen.getByLabelText(`add-new-song-section-part-card-${section.name}`))
    await user.click(screen.getByRole('button', { name: 'close' }))

    expect(screen.getByLabelText(`add-new-song-section-part-card-${section.name}`)).toBeInTheDocument()
    expect(screen.queryByRole('textbox', { name: /name/i })).not.toBeInTheDocument()
  })

  it('should disable the band member select when the artist is not a band', async () => {
    const user = userEvent.setup()

    render(<AddNewSongSectionPart section={section} />, {
      ...emptySong,
      artist: { ...emptyArtist, isBand: false, bandMembers },
      settings: emptySongSettings
    })

    await user.click(screen.getByLabelText(`add-new-song-section-part-card-${section.name}`))

    expect(screen.getByRole('button', { name: 'select-band-member' })).toBeDisabled()
  })

  it('should show default band member and instrument from song settings', async () => {
    const user = userEvent.setup()
    const settings = {
      ...emptySongSettings,
      defaultBandMember: bandMembers[1],
      defaultInstrument: instruments[1]
    }

    render(<AddNewSongSectionPart section={section} />, {
      ...emptySong,
      artist: { ...emptyArtist, isBand: true, bandMembers },
      settings
    })

    await user.click(screen.getByLabelText(`add-new-song-section-part-card-${section.name}`))

    expect(
      screen.getByRole('button', { name: settings.defaultBandMember.name })
    ).toBeInTheDocument()
    expect(
      screen.getByRole('button', { name: settings.defaultInstrument.name })
    ).toBeInTheDocument()
  })

  it('should send the create request when a part is added', async () => {
    const user = userEvent.setup()

    const songId = 'song-1'
    const newName = 'Part 1'

    let capturedRequest: CreateSongPartRequest
    server.use(
      http.post('/songs/parts', async ({ request }) => {
        capturedRequest = (await request.json()) as CreateSongPartRequest
        return HttpResponse.json({ message: 'part added!' })
      })
    )

    render(withToastify(<AddNewSongSectionPart section={section} />), {
      ...emptySong,
      id: songId,
      artist: { ...emptyArtist, isBand: true, bandMembers },
      settings: emptySongSettings
    })

    await user.click(screen.getByLabelText(`add-new-song-section-part-card-${section.name}`))
    await user.type(screen.getByRole('textbox', { name: /name/i }), newName)
    await user.click(screen.getByRole('button', { name: 'Add' }))

    expect(capturedRequest).toStrictEqual({
      name: newName,
      songId,
      sectionId: section.id
    })
    expect(screen.getByText(`${newName} added!`)).toBeInTheDocument()
    expect(screen.getByLabelText(`add-new-song-section-part-card-${section.name}`)).toBeInTheDocument()
  })

  it('should send the create request when a part is added band member and instrument', async () => {
    const user = userEvent.setup()

    const songId = 'song-1'
    const newName = 'Part 1'

    let capturedRequest: CreateSongPartRequest
    server.use(
      http.post('/songs/parts', async ({ request }) => {
        capturedRequest = (await request.json()) as CreateSongPartRequest
        return HttpResponse.json({ message: 'part added!' })
      })
    )

    render(withToastify(<AddNewSongSectionPart section={section} />), {
      ...emptySong,
      id: songId,
      artist: { ...emptyArtist, isBand: true, bandMembers },
      settings: emptySongSettings
    })

    await user.click(screen.getByLabelText(`add-new-song-section-part-card-${section.name}`))

    await user.click(screen.getByRole('button', { name: 'select-band-member' }))
    await user.click(await screen.findByRole('option', { name: bandMembers[0].name }))

    await user.click(screen.getByRole('button', { name: 'select-instrument' }))
    await user.click(await screen.findByRole('option', { name: instruments[0].name }))

    await user.type(screen.getByRole('textbox', { name: /name/i }), newName)
    await user.click(screen.getByRole('button', { name: 'Add' }))

    expect(capturedRequest).toStrictEqual({
      name: newName,
      songId,
      sectionId: section.id,
      bandMemberId: bandMembers[0].id,
      instrumentId: instruments[0].id
    })
  })

  it('should display a validation error when the part name is empty', async () => {
    const user = userEvent.setup()

    render(<AddNewSongSectionPart section={section} />)

    await user.click(screen.getByLabelText(`add-new-song-section-part-card-${section.name}`))
    await user.click(screen.getByRole('button', { name: 'Add' }))

    expect(screen.getByRole('textbox', { name: /name/i })).toBeInvalid()
  })

  it('should clear the validation error when the form is closed and reopened', async () => {
    const user = userEvent.setup()

    render(<AddNewSongSectionPart section={section} />)

    await user.click(screen.getByLabelText(`add-new-song-section-part-card-${section.name}`))
    await user.click(screen.getByRole('button', { name: 'Add' }))
    expect(screen.getByRole('textbox', { name: /name/i })).toBeInvalid()

    await user.click(screen.getByRole('button', { name: 'close' }))
    await user.click(screen.getByLabelText(`add-new-song-section-part-card-${section.name}`))

    expect(screen.getByRole('textbox', { name: /name/i })).not.toBeInvalid()
  })
})
