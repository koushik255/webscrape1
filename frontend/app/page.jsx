"use client";

import { useState, useEffect } from "react";
1
export default function Home() {
  const [players, setPlayers] = useState([]);
  const [player, setPlayer] = useState ([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetch("http://localhost:3000/players")
      .then((response) => {
        if (!response.ok) {
          throw new Error("Network response was not ok");
        }
        return response.json();
      })
      .then((result) => {
        // Store the actual array of players
        setPlayers(Array.isArray(result) ? result : [result]);
        setIsLoading(false);
      })
      .catch((error) => {
        setError("Error fetching data: " + error.message);
        setIsLoading(false);
      });
  }, []);

 




  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-8 bg-gray-50">
      <div className="max-w-4xl w-full p-6 bg-white rounded-lg shadow-md">
        <h1 className="text-3xl font-bold mb-6 text-center text-blue-800">
          Soccer Players
          </h1>
        {isLoading ? (
          <div className="flex justify-center py-8">
            <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
          </div>
        ) : error ? (
          <div className="p-4 bg-red-100 border border-red-400 text-red-700 rounded mb-4">
            <p>{error}</p>
          </div>
        ) : (
          <div>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {players.map((player) => (
                <div
                  key={player.id}
                  className="bg-gradient-to-br from-blue-50 to-indigo-50 p-4 rounded-lg shadow hover:shadow-md transition-shadow duration-200 border border-blue-100"
                >
                  <div className="flex items-center mb-2">
                    <div className="w-12 h-12 bg-blue-600 rounded-full flex items-center justify-center text-white font-bold text-xl mr-3">
                      {player.name.charAt(0)}
                    </div>
                    <div>
         <img
        src={player.photo}
        alt="Player Headshot"
        style={{ width: "200px", height: "auto" }}
      />
    </div>
                    <div>
                      <h2 className="text-xl font-semibold text-gray-800">
                        {player.name}
                      </h2>
                      <p className="text-sm text-gray-500">ID: {player.id}</p>
                    </div>
                  </div>
                  
                  <div className="mt-3 pt-3 border-t border-blue-100">
                    <div className="flex justify-between items-center">
                      <span className="text-gray-700 font-medium">Goals:</span>
                      <span className="bg-blue-600 text-white px-3 py-1 rounded-full font-bold">
                        {player.goals}
                      </span>
                      <span className="text-gray-700 font-medium">Assists:</span>
                      <span className="bg-blue-600 text-white px-3 py-1 rounded-full font-bold">
                        {player.assists}
                      </span>
                    </div>
                  </div>
                  
                  <div className="mt-3 text-xs text-gray-500">
                    <p>Created: {new Date(player.created_at).toLocaleDateString()}</p>
                    <p>Updated: {new Date(player.updated_at).toLocaleDateString()}</p>
                  </div>
                </div>
              ))}
            </div>
            
            {players.length === 0 && (
              <div className="text-center py-8 text-gray-500">
                No players found.
              </div>
            )}
            
            <div className="mt-6 text-center text-gray-500 text-sm">
              Total Players: {players.length}
            </div>
          </div>
        )}
      </div>
    </main>
  );
}
