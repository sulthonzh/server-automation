<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use Symfony\Component\Process\Process;
use Symfony\Component\Process\Exception\ProcessFailedException;

class SetUserStorageLimit extends Controller
{
    public function setQuota(Request $request)
    {
        $username = $request->input('username');
        $limitGB = $request->input('limitGB');
        $goAgentPath = env('GO_AGENT_PATH', './go-agent');

        $process = new Process([$goAgentPath, "setquota", "--username={$username}", "--limitGB={$limitGB}"]);
        $process->run();

        if (!$process->isSuccessful()) {
            return response()->json(['success' => false, 'message' => 'Failed to set user quota.'], 500);
        }

        return response()->json(['success' => true, 'message' => 'User quota set successfully.']);
    }
}
