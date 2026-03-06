param(
    [string]$BaseUrl = "http://localhost:8080",
    [Parameter(Mandatory = $true)]
    [string]$AuthorToken,
    [Parameter(Mandatory = $true)]
    [string]$AdminToken,
    [Parameter(Mandatory = $true)]
    [string]$ProjectId,
    [Parameter(Mandatory = $true)]
    [string]$DocumentId,
    [int]$ChapterNumber = 1
)

$ErrorActionPreference = "Stop"

function Invoke-JsonRequest {
    param(
        [string]$Method,
        [string]$Url,
        [string]$Token,
        [object]$Body
    )

    $headers = @{}
    if ($Token) {
        $headers["Authorization"] = "Bearer $Token"
    }

    if ($null -ne $Body) {
        $json = $Body | ConvertTo-Json -Depth 8
        return Invoke-RestMethod -Method $Method -Uri $Url -Headers $headers -ContentType "application/json" -Body $json
    }

    return Invoke-RestMethod -Method $Method -Uri $Url -Headers $headers
}

function Get-ApiData {
    param([object]$Response)

    if ($null -eq $Response) {
        return $null
    }

    if ($Response.PSObject.Properties.Name -contains "data") {
        return $Response.data
    }

    return $Response
}

Write-Host "1. Author submits project publication for review"
$publishBody = @{
    bookstoreId   = "local"
    categoryId    = "mvp-category"
    tags          = @("mvp", "publication")
    description   = "Publication MVP end-to-end flow"
    coverImage    = ""
    publishType   = "serial"
    freeChapters  = 1
    authorNote    = "Submitted by e2e script"
    enableComment = $true
    enableShare   = $true
}

$projectPublishResp = Invoke-JsonRequest `
    -Method "POST" `
    -Url "$BaseUrl/api/v1/writer/projects/$ProjectId/publish" `
    -Token $AuthorToken `
    -Body $publishBody

$projectRecord = Get-ApiData $projectPublishResp
if ($null -eq $projectRecord -or [string]::IsNullOrWhiteSpace($projectRecord.id)) {
    throw "Project publication record was not returned."
}
Write-Host ("   project record id: {0}, status: {1}" -f $projectRecord.id, $projectRecord.status)

Write-Host "2. Optional: author submits single document publication for review"
$documentPublishBody = @{
    chapterTitle  = "MVP Chapter $ChapterNumber"
    chapterNumber = $ChapterNumber
    isFree        = $true
    authorNote    = "Document publish from e2e script"
}

$documentPublishResp = Invoke-JsonRequest `
    -Method "POST" `
    -Url "$BaseUrl/api/v1/writer/documents/$DocumentId/publish?projectId=$ProjectId" `
    -Token $AuthorToken `
    -Body $documentPublishBody

$documentRecord = Get-ApiData $documentPublishResp
if ($null -eq $documentRecord -or [string]::IsNullOrWhiteSpace($documentRecord.id)) {
    throw "Document publication record was not returned."
}
Write-Host ("   document record id: {0}, status: {1}" -f $documentRecord.id, $documentRecord.status)

Write-Host "3. Admin fetches pending publication queue"
$pendingResp = Invoke-JsonRequest `
    -Method "GET" `
    -Url "$BaseUrl/api/v1/admin/publications/pending?page=1&pageSize=20" `
    -Token $AdminToken `
    -Body $null

$pendingData = Get-ApiData $pendingResp
$pendingItems = @()
if ($pendingData -is [System.Array]) {
    $pendingItems = $pendingData
} elseif ($null -ne $pendingData.items) {
    $pendingItems = $pendingData.items
} elseif ($null -ne $pendingData) {
    $pendingItems = @($pendingData)
}

Write-Host ("   pending count (response page): {0}" -f $pendingItems.Count)

Write-Host "4. Admin approves the project publication"
$reviewBody = @{
    action = "approve"
    note   = "Approved by e2e publication flow script"
}

$reviewResp = Invoke-JsonRequest `
    -Method "POST" `
    -Url "$BaseUrl/api/v1/admin/publications/$($projectRecord.id)/review" `
    -Token $AdminToken `
    -Body $reviewBody

$reviewedRecord = Get-ApiData $reviewResp
if ($null -eq $reviewedRecord) {
    throw "Project review response did not contain a record."
}
Write-Host ("   project review status: {0}" -f $reviewedRecord.status)

Write-Host "5. Read-side verification through bookstore routes"
$bookstoreBookResp = Invoke-JsonRequest `
    -Method "GET" `
    -Url "$BaseUrl/api/v1/bookstore/books/$($reviewedRecord.bookstoreId)" `
    -Token "" `
    -Body $null

$bookstoreBook = Get-ApiData $bookstoreBookResp
if ($null -eq $bookstoreBook) {
    throw "Bookstore book lookup returned no data."
}
Write-Host ("   bookstore book id: {0}, title: {1}" -f $bookstoreBook.id, $bookstoreBook.title)

$chapterListResp = Invoke-JsonRequest `
    -Method "GET" `
    -Url "$BaseUrl/api/v1/bookstore/books/$($bookstoreBook.id)/chapters" `
    -Token "" `
    -Body $null

$chapterList = Get-ApiData $chapterListResp
Write-Host "6. Reader-side chapter access verification"

$readerChapterResp = Invoke-JsonRequest `
    -Method "GET" `
    -Url "$BaseUrl/api/v1/reader/books/$($bookstoreBook.id)/chapters/$($documentRecord.resourceId)" `
    -Token $AuthorToken `
    -Body $null

$readerChapter = Get-ApiData $readerChapterResp

[pscustomobject]@{
    projectPublicationRecordId = $projectRecord.id
    documentPublicationRecordId = $documentRecord.id
    bookstoreBookId = $bookstoreBook.id
    chapterList = $chapterList
    readerChapter = $readerChapter
} | ConvertTo-Json -Depth 10
