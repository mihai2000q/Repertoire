import { ReactElement } from 'react'
import useFixedDocumentTitle from "../hooks/useFixedDocumentTitle.ts";

function Home(): ReactElement {
  useFixedDocumentTitle('Home')

  return (
    <div>
      <h1>Home</h1>
    </div>
  )
}

export default Home
