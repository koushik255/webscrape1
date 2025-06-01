"use client";

import { useState, useEffect } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";

export default function PlayerPage() {
  const { playerName } = useParams();
  const [player, setPlayer] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
  // decode the player name from the URL
  const decodedPlayerName = decodeURIComponent(playerName);
  
  fetch(`http://localhost:3000/players/${decodedPlayerName}`)
    .then((response) => {
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      return response.json();
    })
    .then((result) => {
      // if the API returns an array, take the first player or handle as needed
      setPlayer(Array.isArray(result) && result.length > 0 ? result[0] : result);
      setIsLoading(false);
    })
    .catch(async (error) => {
      // setError("Error fetching data: " + error.message);
      
      try {
        // Ping the 8000 port with the actual player name (not {decodedPlayerName})
        await fetch(`http://localhost:8000/${decodedPlayerName}`, {
          method: 'GET',
          signal: AbortSignal.timeout(5000)
        });
        
        // After successful ping, refresh the page
        window.location.reload();
        
      } catch (pingError) {
        // If ping also fails, just leave the error state
        console.log("Ping failed:", pingError);
        setIsLoading(false);
      }
    });
}, [playerName]);

  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-8 bg-gray-50">
      <div className="max-w-4xl w-full p-6 bg-white rounded-lg shadow-md">
        <Link 
          href="/" 
          className="inline-block mb-6 text-blue-600 hover:text-blue-800 transition-colors"
        >
          ← Back to all players
        </Link>
        
        <h1 className="text-3xl font-bold mb-6 text-center text-blue-800">
          Player Details
        </h1>
        
        {isLoading ? (
          <div className="flex justify-center py-8">
            <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
          </div>
        ) : error ? (
          <div className="p-4 bg-red-100 border border-red-400 text-red-700 rounded mb-4">
            <p>{error}</p>
          </div>
        ) : player ? (
          <div className="bg-gradient-to-br from-blue-50 to-indigo-50 p-6 rounded-lg shadow border border-blue-100">
            <div className="flex flex-col md:flex-row items-center md:items-start gap-6">
              <div className="w-full md:w-1/3 flex justify-center">
                {player.photo ? (
                  <img
                    src={player.photo}
                    alt={`${player.name} Headshot`}
                    className="rounded-lg shadow-md max-w-full h-auto"
                    style={{ maxHeight: "300px" }}
                  />
                ) : (
                  <div className="w-32 h-32 bg-blue-600 rounded-full flex items-center justify-center text-white font-bold text-4xl">
                    {player.name?.charAt(0)}
                  </div>
                )}
              </div>
              
              <div className="w-full md:w-2/3">
                <div className="flex items-center mb-4">
                  <h2 className="text-3xl font-semibold text-gray-800">
                    {player.name}
                  </h2>
                </div>
                
                <div className="bg-white p-4 rounded-lg shadow-sm mb-4">
                  <h3 className="text-lg font-medium text-gray-700 mb-3">Player Stats</h3>
                  <div className="grid grid-cols-2 gap-4">
                    <div className="bg-blue-50 p-4 rounded-lg">
                      <p className="text-gray-600 text-sm">Goals</p>
                      <p className="text-3xl font-bold text-blue-700">{player.goals}</p>
                    </div>
                    <div className="bg-blue-50 p-4 rounded-lg">
                      <p className="text-gray-600 text-sm">Assists</p>
                      <p className="text-3xl font-bold text-blue-700">{player.assists}</p>
                    </div>
                  </div>
                </div>
                
                <div className="bg-white p-4 rounded-lg shadow-sm">
                  <h3 className="text-lg font-medium text-gray-700 mb-3">Player Information</h3>
                  <div className="space-y-2">
                    <div className="flex justify-between">
                      <span className="text-gray-600">Player ID:</span>
                      <span className="font-medium">{player.id}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600">Created:</span>
                      <span className="font-medium">
                        {new Date(player.created_at).toLocaleDateString()}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600">Last Updated:</span>
                      <span className="font-medium">
                        {new Date(player.updated_at).toLocaleDateString()}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        ) : (
          <div className="text-center py-8 text-gray-500">
            Player not found.
          </div>
        )}
      </div>
    </main>
  );
}