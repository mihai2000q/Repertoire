import { emptySongSection, reduxRender } from '../../../../../../../../test-utils.tsx'
import { screen, waitFor } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { expect } from 'vitest'
import { http, HttpResponse } from 'msw'
import { SongSection } from '../../../../../../../../types/models/Song.ts'
import { setupServer } from 'msw/node'
import SongSectionSelect from './SongSectionSelect.tsx'

describe('SongSection Select', () => {
  const songSections: SongSection[] = [
    {
      ...emptySongSection,
      id: '1',
      name: 'Chorus'
    },
    {
      ...emptySongSection,
      id: '2',
      name: 'Verse'
    }
  ]

  const handlers = [
    http.get('/songs/sections', async () => {
      return HttpResponse.json(songSections)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render and change song sections', async () => {
    const user = userEvent.setup()

    const newSongSection = songSections[0]
    const newOption = { label: newSongSection.name, value: newSongSection.id }

    const onChange = vitest.fn()

    const label = 'label'

    const [{ rerender }] = reduxRender(
      <SongSectionSelect label={label} option={null} onOptionChange={onChange} songId={''} />
    )

    const select = screen.getByRole('combobox', { name: label })
    expect(select).toHaveValue('')
    expect(select).toBeDisabled()
    await waitFor(() => expect(select).not.toBeDisabled())
    await user.click(select)

    for (const songSection of songSections) {
      expect(await screen.findByRole('option', { name: songSection.name })).toBeInTheDocument()
    }

    const selectedOption = screen.getByRole('option', { name: newSongSection.name })
    await user.click(selectedOption)

    expect(onChange).toHaveBeenCalledOnce()
    expect(onChange).toHaveBeenCalledWith(newOption)

    rerender(
      <SongSectionSelect label={label} option={newOption} onOptionChange={onChange} songId={''} />
    )

    expect(screen.getByRole('combobox', { name: label })).toHaveValue(newSongSection.name)
  })
})
