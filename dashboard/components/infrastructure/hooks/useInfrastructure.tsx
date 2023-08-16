import { useEffect, useRef, useState } from 'react';
import settingsService from '../../../services/settingsService';
import useToast from '../../toast/hooks/useToast';

export type Vertices = {
  id: string;
  name: string;
  service: string;
  data: any;
};

export type Edges = {
  from: string;
  to: string;
  name: string;
};

function useInfrastructure() {
  const [loading, setLoading] = useState(true);
  const { toast, setToast, dismissToast } = useToast();
  const [inventory, setInventory] = useState<Vertices[]>([]);
  const [edges, setEdges] = useState<Edges[]>([]);

  useEffect(() => {
    let mounted = true;

    settingsService.getAllGlobalResources().then(res => {
      if (res === Error) {
        setToast({
          hasError: true,
          title: `There was an error when fetching the cloud providers`,
          message: `Please refresh the page and try again.`
        });
      } else {
        const allResources = res.map((item: any) => ({
          id: item.id,
          name: item.name,
          service: item.service,
          data: item.data
        }));
        setInventory(allResources);
      }
    });

    settingsService.getAllEdges().then(res => {
      if (res === Error) {
        setToast({
          hasError: true,
          title: `There was an error when fetching the cloud providers`,
          message: `Please refresh the page and try again.`
        });
      } else {
        const allEdges = res.map((item: any) => ({
          from: item.source,
          to: item.dest,
          name: item.name
        }));
        setEdges(allEdges);
      }
    });

    return () => {
      mounted = false;
    };
  }, []);
  return {
    loading,
    inventory,
    toast,
    setToast,
    edges
  };
}

export default useInfrastructure;
