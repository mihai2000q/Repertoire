import {
  emptySongPart,
  emptySongSettings,
  reduxRender
} from '../../../../../../../../test-utils.tsx'
import SongOutlineSettingsButton from './SongOutlineSettingsButton.tsx'
import { screen, waitFor } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { Instrument, SongPart } from '../../../../../../../../types/models/Song.ts'
import { BandMember } from '../../../../../../../../types/models/Artist.ts'
import { UpdateSongSettingsRequest } from '../../../../../../../../types/requests/SongRequests.ts'
import { UpdateAllSongPartsRequest } from '../../../parts/types/requests/SongPartRequests.ts'

describe('Song Parts Settings Button', () => {
  const instruments: Instrument[] = [
    {
      id: '1',
      name: 'Guitar'
    },
    {
      id: '2',
      name: 'Piano'
    },
    {
      id: '3',
      name: 'Flute'
    }
  ]

  const bandMembers: BandMember[] = [
    {
      id: '1',
      name: 'Chester',
      roles: []
    },
    {
      id: '2',
      name: 'Michael',
      roles: []
    }
  ]

  const handlers = [
    http.get('/songs/instruments', () => {
      return HttpResponse.json(instruments)
    }),
    http.put('/songs/settings', () => {
      return HttpResponse.json({ message: 'it worked' })
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    const user = userEvent.setup()

    const [{ rerender }] = reduxRender(
      <SongOutlineSettingsButton settings={emptySongSettings} parts={[]} songId={''} />
    )

    expect(screen.getByRole('button', { name: 'settings' })).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'settings' }))
    expect(await screen.findByRole('dialog')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'select-instrument' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'select-band-member' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'select-band-member' })).toBeDisabled()

    rerender(
      <SongOutlineSettingsButton
        settings={emptySongSettings}
        parts={[]}
        songId={''}
        bandMembers={bandMembers}
      />
    )

    expect(screen.getByRole('button', { name: 'select-band-member' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'select-band-member' })).not.toBeDisabled()
  })

  it('should render and display current settings', async () => {
    const user = userEvent.setup()

    const defaultInstrument = instruments[1]
    const defaultBandMember = bandMembers[1]

    reduxRender(
      <SongOutlineSettingsButton
        settings={{ ...emptySongSettings, defaultInstrument, defaultBandMember }}
        parts={[]}
        songId={''}
        bandMembers={bandMembers}
      />
    )

    await user.click(screen.getByRole('button', { name: 'settings' }))

    expect(screen.getByRole('button', { name: defaultInstrument.name })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: defaultBandMember.name })).toBeInTheDocument()
  })

  describe('band member', () => {
    it('should send update song settings request when changing the band member and display updated popover', async () => {
      const user = userEvent.setup()

      let capturedRequest: UpdateSongSettingsRequest
      server.use(
        http.put('/songs/settings', async (req) => {
          capturedRequest = (await req.request.json()) as UpdateSongSettingsRequest
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      const songSettings = {
        ...emptySongSettings,
        id: 'some-id'
      }
      const newBandMember = bandMembers[1]

      const parts: SongPart[] = [{ ...emptySongPart, bandMember: bandMembers[0] }]

      reduxRender(
        <SongOutlineSettingsButton
          settings={songSettings}
          parts={parts}
          songId={''}
          bandMembers={bandMembers}
        />
      )

      await user.click(screen.getByRole('button', { name: 'settings' }))
      await user.click(screen.getByRole('button', { name: 'select-band-member' }))
      await user.click(await screen.findByRole('option', { name: newBandMember.name }))

      await waitFor(() =>
        expect(capturedRequest).toStrictEqual({
          settingsId: songSettings.id,
          defaultBandMemberId: newBandMember.id
        })
      )

      expect(screen.getByText(/update all parts' band members/i)).toBeInTheDocument()
    })

    it('should send update all song parts request when changing the band member and accepting the band members popover', async () => {
      const user = userEvent.setup()

      let capturedRequest: UpdateAllSongPartsRequest
      server.use(
        http.put('/songs/parts/all', async (req) => {
          capturedRequest = (await req.request.json()) as UpdateAllSongPartsRequest
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      const newBandMember = bandMembers[1]
      const songId = 'some-song-id'

      const parts: SongPart[] = [{ ...emptySongPart, bandMember: bandMembers[0] }]

      reduxRender(
        <SongOutlineSettingsButton
          settings={emptySongSettings}
          parts={parts}
          songId={songId}
          bandMembers={bandMembers}
        />
      )

      await user.click(screen.getByRole('button', { name: 'settings' }))
      await user.click(screen.getByRole('button', { name: 'select-band-member' }))
      await user.click(await screen.findByRole('option', { name: newBandMember.name }))

      await user.click(await screen.findByRole('button', { name: 'confirm' }))
      await waitFor(() =>
        expect(capturedRequest).toStrictEqual({
          songId: songId,
          bandMemberId: newBandMember.id
        })
      )
    })

    it('should not always display the updated band member popover (when there are no parts with distinct members)', async () => {
      const user = userEvent.setup()

      const newBandMember = bandMembers[1]
      const songId = 'some-song-id'

      const parts: SongPart[] = [{ ...emptySongPart, bandMember: newBandMember }]

      reduxRender(
        <SongOutlineSettingsButton
          settings={emptySongSettings}
          parts={parts}
          songId={songId}
          bandMembers={bandMembers}
        />
      )

      await user.click(screen.getByRole('button', { name: 'settings' }))
      await user.click(screen.getByRole('button', { name: 'select-band-member' }))
      await user.click(await screen.findByRole('option', { name: newBandMember.name }))

      expect(screen.queryByText(/update all parts' band members/i)).not.toBeInTheDocument()
    })
  })

  describe('instrument', () => {
    it('should send update song settings request when changing the instrument and display updated popover', async () => {
      const user = userEvent.setup()

      let capturedRequest: UpdateSongSettingsRequest
      server.use(
        http.put('/songs/settings', async (req) => {
          capturedRequest = (await req.request.json()) as UpdateSongSettingsRequest
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      const songSettings = {
        ...emptySongSettings,
        id: 'some-id'
      }
      const newInstrument = instruments[1]

      const parts: SongPart[] = [{ ...emptySongPart, instrument: instruments[0] }]

      reduxRender(
        <SongOutlineSettingsButton settings={songSettings} parts={parts} songId={''} />
      )

      await user.click(screen.getByRole('button', { name: 'settings' }))
      await user.click(screen.getByRole('button', { name: 'select-instrument' }))
      await user.click(await screen.findByRole('option', { name: newInstrument.name }))

      await waitFor(() =>
        expect(capturedRequest).toStrictEqual({
          settingsId: songSettings.id,
          defaultInstrumentId: newInstrument.id
        })
      )

      expect(screen.getByText(/update all parts' instruments/i)).toBeInTheDocument()
    })

    it('should send update all song parts request when changing the instrument and accepting the instruments popover', async () => {
      const user = userEvent.setup()

      let capturedRequest: UpdateAllSongPartsRequest
      server.use(
        http.put('/songs/parts/all', async (req) => {
          capturedRequest = (await req.request.json()) as UpdateAllSongPartsRequest
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      const newInstrument = instruments[1]
      const songId = 'some-song-id'

      const parts: SongPart[] = [{ ...emptySongPart, instrument: instruments[0] }]

      reduxRender(
        <SongOutlineSettingsButton
          settings={emptySongSettings}
          parts={parts}
          songId={songId}
        />
      )

      await user.click(screen.getByRole('button', { name: 'settings' }))
      await user.click(screen.getByRole('button', { name: 'select-instrument' }))
      await user.click(await screen.findByRole('option', { name: newInstrument.name }))

      await user.click(await screen.findByRole('button', { name: 'confirm' }))
      await waitFor(() =>
        expect(capturedRequest).toStrictEqual({
          songId: songId,
          instrumentId: newInstrument.id
        })
      )
    })

    it('should not always display the updated instrument popover (when there are no parts with distinct instruments)', async () => {
      const user = userEvent.setup()

      const newInstrument = instruments[1]
      const songId = 'some-song-id'

      const parts: SongPart[] = [{ ...emptySongPart, instrument: newInstrument }]

      reduxRender(
        <SongOutlineSettingsButton
          settings={emptySongSettings}
          parts={parts}
          songId={songId}
        />
      )

      await user.click(screen.getByRole('button', { name: 'settings' }))
      await user.click(screen.getByRole('button', { name: 'select-instrument' }))
      await user.click(await screen.findByRole('option', { name: newInstrument.name }))

      expect(screen.queryByText(/update all parts' instruments/i)).not.toBeInTheDocument()
    })
  })
})
