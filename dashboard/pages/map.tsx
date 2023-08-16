import Head from 'next/head';
import { ReactNode, useContext, useState } from 'react';
import InfrastructureLayout from '../components/infrastructure/components/InfrastructureLayout';
import InfrastructureGraph from '../components/infrastructure/components/InfrastructureGraph';
import useInfrastructure from '../components/infrastructure/hooks/useInfrastructure';

export default function Infrastructure() {
  const { loading, inventory, toast, setToast, edges } = useInfrastructure();
  return (
    <div>
      <Head>
        <title>Infrastructure - Komiser</title>
        <meta name="description" content="Infrastructure - Komiser" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <InfrastructureLayout>
        {inventory && inventory.length > 0 && (
          <InfrastructureGraph vertices={inventory} edges={edges} />
        )}
      </InfrastructureLayout>
    </div>
  );
}
